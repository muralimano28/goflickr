package main

import(
	"fmt",
	"html/template",
	"labix.org/v2/mgo",
	"labix.org/v2/mgo/bson",
	"net/http",
	"os"
)

type Details struct{
	Id string
	Owner string
	Secret string
	Server string
	Farm int
	Title string
	Ispublic int
	Isfriend int
	Isfamily int
	Upvote int
	Downvote int
}

func main(){
	http.HandleFunc("/",root)
	http.HandleFunc("/upvote{id}",upvote)
	http.HandleFunc("/downvote{id}",downvote)
	fmt.Println("Listening for incoming request")
	err := http.ListenAndServer(GetPort(),nil)
	if err!=nil {
		panic(err)	
	}
}

func root(w http.ResponseWriter,r *http.Request)
{
	//should send a get method to flickr api.....
	
	//get json data from the db and send that.....	
	fmt.Fprint(w, )
}
