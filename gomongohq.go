package main

import(
	"fmt",
	"html/template",
	"labix.org/v2/mgo",
	"labix.org/v2/mgo/bson",
	"net/http",
	"os",
	"github.com/unrolled/render"
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
	fmt.Println("Listening for incoming request")
	err := http.ListenAndServer(GetPort(),nil)
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


func root(w http.ResponseWriter,r *http.Request)
{
	//..............Sending a get method to flickr api......................

	response, err := http.Get("https://api.flickr.com/services/rest/?method=flickr.photos.search&api_key=c2ef6776bc7946440b01cc9070d55ac0&tags=cute+puppy&per_page=1&format=json&nojsoncallback=1&auth_token=72157654809259522-833770d09d600311&api_sig=164482f059ef2a032b504551e720b66c")
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
		fmt.Println(arrayofjsondata);
   	}

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
	for i,v := range arrayofjsondata {
		doc := Details{Id:v.Id, Owner:v.Owner, Secret:v.Secret, Server:v.Server, Farm:v.Farm, Title:v.Title, Ispublic:v.Ispublic, Isfriend:v.Isfriend, Isfamily:v.Isfamily, Upvote:0, Downvote:0}
		err := collection.Insert(doc)
		if err != nil {
			fmt.Printf("Cant insert into document : %v",err)
			os.Exit(1)		
		}		
	}

	//...........................Should get data from the DB................................

	result := Details{}
	errn := collection.Find(bson.M{}).One(&result)
	if errn != nil {
		fmt.Printf("Cant get data using Find. The error is : %v",errn)
		os.Exit(1)	
	}

	//............................showing data in console..................................

	fmt.Printf("%+v\n",result)

	//............................sending data to the client.............................

	
	r.JSON(w,http.StatusOK,&result)
	
}

