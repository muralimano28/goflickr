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
	"strconv"
    "strings"
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
	Time int64
}

func main(){

	//.......................... Routing each request...........................

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

	//..............Setting tag value in flickr uri...............

	usrtag := r.URL.Query().Get("tag")
    output := strings.Replace(usrtag," ","+",-1)
	fmt.Println("The user tag is ",output)
	var seedtime int64
	var pgno int
	if len(output)==0 {
		output = "cute+puppy"
		seedtime = time.Now().UnixNano() / int64(time.Millisecond)
		rand.Seed(seedtime)	
		pgno = rand.Intn(143)
	}else{
		pgno = 1
	}

	
	//.................creating a random number to fetch random page from flickr......................

	pgstring := strconv.Itoa(pgno)
	fmt.Println("The pgno is ",pgstring) 
	uri := "https://api.flickr.com/services/rest/?method=flickr.photos.search&api_key=c3e4b4f1288dd22b93d0eb607feac333&tags="+output+"&page="+pgstring+"&per_page=50&format=json&nojsoncallback=1"
	fmt.Println(uri)

	//................Getting data from flickr using uri..........................

	response, err := http.Get(uri)
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
		fmt.Println("arrayofjsondata is ",arrayofjsondata)

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
	
		//...........................If the document count goes more than 1500000 it will not store in database.....................
		//...........................bcoz am using free db from mongolab....limitation..:-P.........................................

		//count,err := collection.Count()
		//if count==1500000{	
			/*for i,v := range arrayofjsondata {
				time := time.Now().UnixNano() / int64(time.Millisecond)
				doc := PuppyDetails{Id:v.Id, Owner:v.Owner, Secret:v.Secret, Server:v.Server, Farm:v.Farm, Title:v.Title, Ispublic:v.Ispublic, Isfriend:v.Isfriend, Isfamily:v.Isfamily, Upvote:0, Downvote:0, Time: time}
				if i<5 {
					fmt.Println("i --> v --> ",i,v,doc)
				}
				err := collection.Insert(doc)
				if err != nil {
					fmt.Printf("Cant insert into document : %v",err)
					os.Exit(1)		
				}		
			}*/
			for i := 0; i<len(arrayofjsondata); i++ {
				time := time.Now().UnixNano() / int64(time.Millisecond)
				doc := PuppyDetails{}
				fmt.Println("initally the doc is ",doc)
				doc = PuppyDetails{Id:arrayofjsondata[i].Id, Owner:arrayofjsondata[i].Owner, Secret:arrayofjsondata[i].Secret, Server:arrayofjsondata[i].Server, Farm:arrayofjsondata[i].Farm, Title:arrayofjsondata[i].Title, Ispublic:arrayofjsondata[i].Ispublic, Isfriend:arrayofjsondata[i].Isfriend, Isfamily:arrayofjsondata[i].Isfamily, Upvote:0, Downvote:0, Time: time}
				if i<5 {
					fmt.Println("i ",i,"--> ",doc)
				}
				err := collection.Insert(doc)
				if err != nil {
					fmt.Printf("Cant insert into document : %v",err)
					os.Exit(1)		
				}			
			}
		//}

		//...........................fetching data from the DB................................

		result := []PuppyDetails{}
		errm := collection.Find(bson.M{}).Sort("-time").Limit(50).All(&result)
		if errm != nil {
			fmt.Printf("Cant get data using Find. The error is : %v",errm)
			os.Exit(1)	
		}

		//............................sending data in JSON format to the client.............................
	
		ren := render.New(render.Options{
	        	IndentJSON: true,
	    	})
		ren.JSON(w,http.StatusOK,&result)
	}	
}

func upvote(w http.ResponseWriter, r *http.Request){

	//..............................getting id from the request uri........................................

	usrid := r.URL.Query().Get("id")
	if len(usrid)==0 {
		fmt.Println("user id is not recieved")
		os.Exit(1)	
	}
	
	//..............connecting to the db........................

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
	_,errn := collection.UpdateAll(bson.M{"id":usrid},bson.M{"$inc":bson.M{"upvote":1}})
	if errn != nil {
		fmt.Printf("Error while updating for upvote : %v",errn)	
	}
}

func downvote(w http.ResponseWriter, r *http.Request){

	//.........................Getting id from the request uri..............................

	usrid := r.URL.Query().Get("id")
	if len(usrid)==0 {
		fmt.Println("user id is not recieved")
		os.Exit(1)	
	}
	
	//................Connecting to the db........................

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
	_,errn := collection.UpdateAll(bson.M{"id":usrid},bson.M{"$inc":bson.M{"downvote":1}})
	if errn != nil {
		fmt.Printf("Error while updating for upvote : %v",errn)	
	}

}

