package main

import(
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"os"
	"io/ioutil"
	"github.com/unrolled/render"
	"encoding/json"
	"math/rand"
	"time"
)

type Details struct{
	Id string	`json:"id"`
	Owner string	`json:"owner"`
	Secret string	`json:"secret"`
	Server string	`json:"server"`
	Farm int	`json:"farm"`
	Title string 	`json:"title"`
	Ispublic int	`json:"ispublic"`
	Isfriend int	`json:"isfriend"`
	Isfamily int	`json:"isfamily"`
	Time int64	`json:"time"`
}
type Packet struct{
	Page int	`json:"page"`
	Pages int	`json:"pages"`
	Perpage int	`json:"perpage"`
	Total string	`json:"total"`
	Photo []Details	`json:"photo"`
} 
type Message struct{
	Photos Packet	`json:"photos"`
	Stat string	`json:"stat"`
}

type PuppyDetails struct{
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
	http.HandleFunc("/fetch",root)
	http.HandleFunc("/upvote",upvote)
	http.HandleFunc("/downvote",downvote)
	fmt.Println("Listening for incoming request")
	err := http.ListenAndServe(GetPort(),nil)
	if err!=nil {
		panic(err)	
	}
}

//...............................Getting port for Heroku.......................

func GetPort() string{
	var port = os.Getenv("PORT")
	if port=="" {
		port = "4747"
		fmt.Println("INFO: no PORT environment variable is available")	
	}
	return ":" + port
}

func root(w http.ResponseWriter,r *http.Request){
	//..............Sending a get method to flickr api......................
	
	//..............get the value in url...............

	usrtag := r.URL.Query().Get("tag")
	if len(usrtag)==0 {
		usrtag = "cute+puppy"
	}
	
	//.................creating a random number for page......................
	rand.Seed(time.Now().UnixNano() / int64(time.Millisecond))	
	pgno := rand.Intn(7114) 
	fmt.Println("The page number is",pgno)
	uri := "https://api.flickr.com/services/rest/?method=flickr.photos.search&api_key=c3e4b4f1288dd22b93d0eb607feac333&tags="+usrtag+"&page="+string(pgno)+"&per_page=100&format=json&nojsoncallback=1"

	//fmt.Println("The random generated is ",pgno,"and the uri is ",uri)
	response, err := http.Get(uri)
	//response, err := http.Get("https://api.flickr.com/services/rest/?method=flickr.photos.search&api_key=6fcd9a3cf8a8a260b36df54274e6b5bc&tags=cute+puppy&page=2&format=json&nojsoncallback=1&auth_token=72157654440984338-11b28d428a9e9b02&api_sig=e434c73b0d14e0790d328d63e9153052")
    	if err != nil {
        	fmt.Printf("%s", err)
        	os.Exit(1)
    	} else {
        	defer response.Body.Close()
        	contents, err := ioutil.ReadAll(response.Body)
        	if err != nil {
            		fmt.Printf("%s", err)
            		os.Exit(1)
        	}
		//.....................Unmarshalling the message..............................
		m := Message{}
		errn := json.Unmarshal([]byte(contents),&m)
		if errn != nil {
			fmt.Println(errn);
			os.Exit(1)
		}
		//..........arrayofjsondata holds the data that we need to upload...................
		arrayofjsondata := m.Photos.Photo;
		//fmt.Println(arrayofjsondata);
   	//}

	//..............Now we need to connect to the db........................

	uri := os.Getenv("MONGOHQ_URL")
	if uri == "" {
		fmt.Println("No connection string provided");
		os.Exit(1)
	}	
	
	sess, err := mgo.Dial(uri)
	if err != nil {
		fmt.Println("Cant connect to Mongodb using the given uri");
		os.Exit(1)	
	}
	defer sess.Close()

	sess.SetSafe(&mgo.Safe{})

	//...........................Inserting document into the collection in DB.................

	collection := sess.DB("puppypull").C("puppydetails")

	//count,err := collection.Count()
	//if count==0{
		for _,v := range arrayofjsondata {
			doc := PuppyDetails{Id:v.Id, Owner:v.Owner, Secret:v.Secret, Server:v.Server, Farm:v.Farm, Title:v.Title, Ispublic:v.Ispublic, Isfriend:v.Isfriend, Isfamily:v.Isfamily, Upvote:0, Downvote:0}
			err := collection.Insert(doc)
			if err != nil {
				fmt.Printf("Cant insert into document : %v",err)
				os.Exit(1)		
			}		
		}
	//}
	//...........................Should get data from the DB................................

	result := []PuppyDetails{}
	errm := collection.Find(bson.M{}).Sort("-$oid").Limit(100).All(&result)
	if errm != nil {
		fmt.Printf("Cant get data using Find. The error is : %v",errm)
		os.Exit(1)	
	}

	//............................showing data in console..................................

	//fmt.Printf("%+v\n",result)

	//............................sending data to the client.............................

	ren := render.New(render.Options{
        	IndentJSON: true,
    	})
	ren.JSON(w,http.StatusOK,&result)
	}	
}

func upvote(w http.ResponseWriter, r *http.Request){
	usrid := r.URL.Query().Get("id")
	if len(usrid)==0 {
		fmt.Println("user id is not recieved")
		os.Exit(1)	
	}
	//fmt.Println("This is the id received")
	//..............Now we need to connect to the db........................

	uri := os.Getenv("MONGOHQ_URL")
	if uri == "" {
		fmt.Println("No connection string provided");
		os.Exit(1)
	}	
	
	sess, err := mgo.Dial(uri)
	if err != nil {
		fmt.Println("Cant connect to Mongodb using the given uri");
		os.Exit(1)	
	}
	defer sess.Close()

	sess.SetSafe(&mgo.Safe{})

	//................Now updating the document in collection.......................

	collection := sess.DB("puppypull").C("puppydetails")
	errn := collection.Update(bson.M{"id":usrid},bson.M{"$inc":bson.M{"upvote":1}})
	if errn != nil {
		fmt.Printf("Error while updating for upvote : %v",errn)	
	}

}

func downvote(w http.ResponseWriter, r *http.Request){
	usrid := r.URL.Query().Get("id")
	if len(usrid)==0 {
		fmt.Println("user id is not recieved")
		os.Exit(1)	
	}
	fmt.Println("This is the id received")
	//..............Now we need to connect to the db........................

	uri := os.Getenv("MONGOHQ_URL")
	if uri == "" {
		fmt.Println("No connection string provided");
		os.Exit(1)
	}	
	
	sess, err := mgo.Dial(uri)
	if err != nil {
		fmt.Println("Cant connect to Mongodb using the given uri");
		os.Exit(1)	
	}
	defer sess.Close()

	sess.SetSafe(&mgo.Safe{})

	//................Now updating the document in collection.......................

	collection := sess.DB("puppypull").C("puppydetails")
	errn := collection.Update(bson.M{"id":usrid},bson.M{"$inc":bson.M{"downvote":1}})
	if errn != nil {
		fmt.Printf("Error while updating for upvote : %v",errn)	
	}

}

