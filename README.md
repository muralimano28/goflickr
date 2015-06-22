# Go Server
Introduction
Puppy pull is a simple web application that fetches photos of cute puppies from the internet and displays it on a webpage. The visitors to the page can upvote and downvote an image. They can sort the images according to the most voted. Every image will display the number of votes it has received.Go server helps puppy pull by fetching data from flickr and giving it to Nodejs and then to client for display.

Highlights :
Client   ---> Node Server ---> Go server ---> MongoDB (Puppypull) 

Goserver Fetched data form flickr and stores it in MongoDB (Created a free db in Mongolab.com)

Client sends a request to Node server and Node server inturn connects to Go server. Go server fetches data from internet using Flickr API and stores in Mongo database. Go server serves back Node server with JSON data. Node server handles the incoming JSON data and combines it with HTML ,Angular JS and displays it to the client.

To make it work :

1. Clone this respository to your local machine. (Machine should have go installed)
2. Run the command "go run gomongohq.go" in terminal.
3. Now open browser and go to the given url. 
    http://localhost:4747/fetch
This will give a json output.
4. routing for upvote and downvote also handled in gomongohq.go file.

Puppy pull code is in the below repository.

https://github.com/muralimano28/puppypull/
