Implement a model application for customer support ticketing system:

There are incoming requests that are processed by the workers (customer
support agents).

Each request has the following fields:
 - question (string - the text of the customer's question)
 - priority (int - 1, 2, 3)

The tickets with priority 1 are the most important ones, the tickets with
 priority = 3 are the least important ones.

There should be a pool of workers that process incoming questions (the number
of workers should be configurable), and return the reply back to the customer.
Processing takes random time (1-3 seconds), the reply can be filled with a
random string.

If a request is waiting for the response longer than a configured timeout value,
  error should be returned.

If you have time, eventually it can be a RESTful API service with methods:
 - /api/post - posting a question (POST request, accepts json:
    {
      "priority": <priority>,
      "question": <text of the question>
    }
    And returns the reply as json:
    {
      "success": <true|false>,
      "message": <reply on success|error if waiting time exceeded timeout>
    }

-  /api/status - returns the status of the system as a json:
    {
      "workers": <number of active workers">,
      "questions_processed": <total number of processed questions>,
      "average_response_time": <average time for a request to be processed>,
      "queue_length":
        {
          1: <length of first priority queue>,
          2: <length of second priority queue>,
          1: <length of third priority queue>
        }
    }

But if you don't have time, you can just create code for this - for example,
if your main app struct is called App, you can create methods:
  - App.Post(priority int, question string) (struct, error) - error should
    be returned if timeout exceeded

  - App.Status() struct


Think about which structures to use, how to structure code, how to implement
prioritised queues, how to implement timeout.

