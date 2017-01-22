//updated
var request = require('request');
var questionJSON = {};

var code, result, poll;
exports.handler = (event, context) => 
{
	try
	{	
		if(event.session.new)
			{ console.log("New Session Created") }

		switch(event.request.type)
		{
			case "LaunchRequest" : 
				var wittyResponse;
				var randomint =  Math.floor(Math.random()*6+1);

				switch(randomint)
				{
					case 1,2,3 : 
						var wittyResponse = `I like your hair by the way.`;
						break;

					case 3,4,5 : 
						var wittyResponse = `Your shirt really brings out your eyes.`;
						break;

					default : 
						var wittyResponse = `It is a nice day outside! We should go for a walk. Wait I do not have legs. `;
				}
				console.log("Launch Requested");
				context.succeed(
					generateResponse(
						buildSpeechResponse(`Hello I am Kalos. ` + wittyResponse, true)),
						{} 
					);
			break;

			case "IntentRequest" : 
					switch(event.request.intent.name)
					{
						case "CreatePoll" : 
							createPoll();
							context.succeed(
								generateResponse(
									buildSpeechResponse(`Your poll has been created. What is your question you would like to add to your new poll?`, true),
									{}
								)
							)
						break;

						case "SetQuestion" : 
							console.log("Adding Question to Poll");
							var question = event.request.intent.slots.question;
							console.log(question);
							setQuestion(question); 

							
								context.succeed(
									generateResponse(
										buildSpeechResponse(`Your question was added to the poll! What are the options you would like to add to the question?`, true),
										{}
									)
								)
						break;

						case "AddOption" : 
							console.log("Adding Options to Poll");
							var option = event.request.intent.slots.option;
							var output = `Your option. `+ option.value +` Was added to the poll! What are the options you would like to add to the question?`;

							if (option.value ==  "Dark Knight")
							{
								output = `Hmm. Not the wisest choice my child.... But `+ option.value +` was added to the poll. What are more options you want to add?`;
							}
							else
							if (option.value ==  "Doctor Strange")
							{
								output = `Ahhhh, take me to the theater too! Anyways, your option `+ option.value +` was added to the poll. What are more options you want to add?`;
							}
							else
							if(option.value ==  "The Avengers")
							{
								output = `Probably the most overrated movie that I have seen! Anyways, your option `+ option.value +` was added to the poll. What are more options you want to add?`;
							}

							setOption(option);
								context.succeed(
									generateResponse(
										buildSpeechResponse(output, true),
										{}
									)
								)
						break;

						case "ReviewPoll" :
							context.succeed(
									generateResponse(
										buildSpeechResponse(`Here is what your poll looks like so far. ` + poll, true),
										{}
									)
								)
						break;

						case "SubmitPoll" :
							request.post(
							    'http://a2ba0d05.ngrok.io:80/createPoll',
							    { json: poll },
							    function (error, response, body) {
							        if (!error && response.statusCode == 200) {
							            console.log(body);
							            code = body;
							        }
							        context.succeed(
										generateResponse(
											buildSpeechResponse(`Successfully submitted your poll! Your poll unique identifier is ` + code, false),
											{}
										)
									)
							    }
							)
						break;

						case "ViewResults" :
							var pollcode = event.request.intent.slots.pollcode;

							console.log(pollcode);
							request.post(
							    'http://a2ba0d05.ngrok.io:80/viewResults',
							    { json: pollcode },
							    function (error, response, body) {
							        if (!error && response.statusCode == 200) {
							            console.log(body);
							            result = body;
							        }
							        else
							        {
							   			 context.succeed(
												generateResponse(
													buildSpeechResponse(`No poll found. Maybe try another code.`, true),
													{}
												)
											)
							   			 return;
							        }
							        var output = `The results for poll `+ pollcode.value + ` has ` + result.responses.length + ` responses. ` + ` The first choice is ` + result.results[0].option + ` and has a score of ` + result.results[0].rating + `. The second choice is ` + result.results[1].option + ` and has a score of ` + result.results[1].rating; 
							        context.succeed(
										generateResponse(
											buildSpeechResponse(output, true),
											{}
										)
									)
									poll = null;
									options = null;
									result = null;
									code = null;
							    }
							)

						break;


						case "DeleteEverything" : 
							poll = null;
							options = null;
							result = null;
							code = null;
							context.succeed(
								generateResponse(
									buildSpeechResponse(`Everything has been erased. So sad.`, true),
									{}
								)
							)
						break;
					};
			break;

			case "SessionEndedRequest":
				console.log(`Session Ended Request`);
				context.succeed(
					generateResponse(
						buildSpeechLetResponse("Byeeeee", true),
						{}
						)
					)
			break;
			break;

			default:
				context.fail(`Invaild Request Name Type: ${event.request.type}`);

		}
	}
	catch (error){ context.fail(`Exception: ${error}`) } 
}

//Outputting
buildSpeechResponse = (outputText, shouldEndSession) =>
{
	return {
		outputSpeech: 
		{
			type: "PlainText",
			text: outputText
		},
		shouldEndSession: shouldEndSession
	}
}

generateResponse = (speechletResponse, sessionAttributes) =>
{
	return {
		version: "1.0",
		sessionAttributes: sessionAttributes,
    	response: speechletResponse
	}
}

//Create Poll
function createPoll()
{
	poll = [];
}

function setQuestion(question)
{
	if(poll == null)
	{
		error(0, "Poll Not Found. Can't set question", true);
		poll = [];
	}
	if(poll.length >= 1)
	{
		error(1, "Too many questions for one poll", false);
		return;
	}

	poll.push(question);
}
function error(errorCode, reason, isFixed)
{
	switch(errorCode)
	{
		case 0 :
		console.log("Not found - " + reason + " - is Fixed: "  + isFixed);
		break;
	}
}

function setOption(option)
{
	if(poll == null)
	{
		error(0, "Poll Not Found. Can't set option", false); return;
	}
	poll.push(option);
}



	/*
	var dataJSON = data;
	var dataSendReady = JSON.stringify(dataJSON);



	var http = require("http");
	var options = {
	host: 'cbb9bd1f.ngrok.io',
	port: 80,
	path: location,
	method: subMethod,
	headers: {
	'Content-Type': 'application/json',
	}
	};
	var request = http.request(options, 
		function(response) 
		{
			console.log('Status: ' + response.statusCode);
			console.log('Headers: ' + JSON.stringify(response.headers));
			response.setEncoding('utf8');
			response.on('data', 
				function (body) 
				{
					console.log('Body: ' + body);
				});
		});

	request.on('error', 
		function(e) {
			console.log('problem with request: ' + e.message);
		});
	// write data to request body
	request.write(dataSendReady);
	request.end();
	*/









