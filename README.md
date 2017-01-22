##Inspiration

When you’re in a group, it’s nearly impossible to decide where to go for lunch or what movie to watch. We've all had experiences where we spent more time browsing through Netflix than actually watching something. The problem is in the way we decide. Asking people to vote once for their favorite option doesn't generate results that reflect the actual preferences of the group. To solve this problem, we created Kalos, a web service that helps groups make smarter decisions quickly and easily.

##How It Works

Kalos makes it easy to create a poll: simply type in your question and add as many options as you want. Once you finish, we generate a code that you can share in your group chat, so others can take your poll. To take a poll, type in your name and the code that was given to you, which brings you to the poll. Using the slider, you can rate as many or as few options as you wish, then hit submit.

##How We Built It

The website runs on Golang. Calculating the results by running the data through our algorithm was done on the server-side. In addition, we integrated Amazon's Alexa with our web service, requiring us to make heavy use of AWS, Lambda, and the Alexa Skills Kit. When the user interacts with Alexa, by detecting key words and phrases, we were able to map the user's request and make an HTTP request to the server.

##Algorithm

The algorithm we developed is based on the alternative vote, but with some enhancements. After each round, the biggest loser is eliminated. The points allocated to that option are discarded if the rating is negative or zero. If the rating is positive, those points are distributed evenly among that user's next highest-rated options. This cycle repeats until there are only two options left.

##What's Next

Although we were able to use AWS and the Alexa Skills Kit to process text input and perform the actions associated with the requests, we ran into issues when we tried to incorporate the Amazon Echo. Unfortunately, the Echo was unable to establish a connection to the Internet, and because we don't have the luxury of unlimited data, our plans to integrate our web service with the Echo were impeded. But given that it was our first time programming Alexa and utilizing AWS, we are very proud of what we were able to accomplish over the course of this hackathon.

##Conclusion

Under current voting systems, rarely do we see results that most people are happy with. Using our advanced algorithm, the information from the same group of users generates results that are much more reflective of the group's preferences. Kalos makes it easy to create and take polls, but most importantly, it generates better results that make more people happy.
