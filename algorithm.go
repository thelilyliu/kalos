package main

import (
	"log"
	"math/rand"
)

/*
  ========================================
  Calculate
  ========================================
*/

func calculateResultsDB(options []string, responses []Response, results []Result) []Result {
	/*
	  ========================================
	  Step 1: Tally up the total points for each option.
	  ========================================
	*/

	for i, option := range options { // loop over each option
		result := new(Result) // create new result
		result.Option = option

		for _, response := range responses { // loop over each response
			result.Rating += response.Ratings[i] // cumulate points
		}

		results = append(results, *result) // append result
	}

	/*
	  ========================================
	  Step 2: Determine if necessary to go through process.
	  ========================================
	*/

	for len(results) > 2 { // while more than two results
		/*
		  ========================================
		  Step 3: Determine the biggest loser.
		  ========================================
		*/

		var minValue float64         // save lowest value, explicitly define type float64
		minValue = results[0].Rating // initialize to first value by default
		var minIndex int             // save index of lowest value

		for i, result := range results { // loop over each result
			if result.Rating < minValue { // if smaller than min
				// set new min
				minValue = result.Rating
				minIndex = i
			}
		}

		/*
		  ========================================
		  Step 4: Discard negative points and save positive points.
		  ========================================
		*/

		points := make([]float64, len(responses))
		// slice to hold distributable points for each response
		// with positive rating of least popular result

		for i, response := range responses { // loop over each response
			if response.Ratings[minIndex] > 0 { // if positive rating
				points[i] = response.Ratings[minIndex] // save points
			}
			// else, negative rating, discard points
		}

		/*
		  ========================================
		  Step 5: Discard option from response and results.
		  ========================================
		*/

		for i, response := range responses { // loop over each response
			// delete element: a = append(a[:i], a[i + 1:]...)
			// delete rating from each response
			response.Ratings = append(response.Ratings[:minIndex], response.Ratings[minIndex+1:]...)
			responses[i] = response // update responses
		}

		// delete rating from results
		results = append(results[:minIndex], results[minIndex+1:]...)

		/*
		  ========================================
		  Step 6: Distribute positive points.
		  ========================================
		*/

		for i, response := range responses { // loop over each response
			if points[i] > 0 { // has points to distribute
				/*
				  ========================================
				  Step 7: Find max value.
				  ========================================
				*/

				var maxValue float64 // save highest value, explicitly define type float64
				maxValue = -2        // initialize to lowest possible value

				for _, rating := range response.Ratings { // loop over each rating
					if rating > maxValue { // if larger than max
						maxValue = rating // set new max
					}
				}

				/*
				  ========================================
				  Step 8: Get number of values that equal max.
				  ========================================
				*/

				counter := 0 // count number of values that equal max

				for _, rating := range response.Ratings { // loop over each rating
					if rating == maxValue { // if value equals max
						counter++ // increase counter
					}
				}

				/*
				  ========================================
				  Step 9: Distribute points among highest ratings.
				  ========================================
				*/

				distribution := points[i] / float64(counter) // number of distributable points per highest rating

				for j, rating := range response.Ratings { // loop over each rating
					if rating == maxValue { // if value equals max
						rating += distribution              // distribute points
						response.Ratings[j] += distribution // update responses
						results[j].Rating += distribution   // update results
					}
				}
			}
		}
	}

	/*
	  ========================================
	  Step 10: Ensure first choice has highest rating.
	  ========================================
	*/

	if results[0].Rating < results[1].Rating { // first choice has lower rating
		results[0], results[1] = results[1], results[0] // swap elements
	}

	return results
}

/*
  ========================================
  Generate Responses
  ========================================
*/

func generateResponsesDB(poll *Poll) {
	number := 4 + rand.Intn(7)
	poll.Responses = make([]Response, number)
	options := len(poll.Options)

	for i := 0; i < number; i++ {
		poll.Responses[i].Name = "a"
		poll.Responses[i].Ratings = make([]float64, options)

		for j := 0; j < options; j++ {
			rating := -2 + rand.Intn(5)
			poll.Responses[i].Ratings[j] = float64(rating)
		}

		response := poll.Responses[i]

		if err := submitResponseDB(poll.ID, &response); err != nil {
			log.Println("generate responses error")
		}
	}
}
