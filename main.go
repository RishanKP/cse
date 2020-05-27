package main

import("fmt"
       "log"
       "os"
       "io/ioutil"
       "strings"
       "go.mongodb.org/mongo-driver/bson"
       "go.mongodb.org/mongo-driver/mongo"
       "go.mongodb.org/mongo-driver/mongo/options"
       "github.com/gorilla/mux"
       "context"
       "net/http"
       "html/template")
type Blog struct{
  Name string
  Subject string
  Sem string
  Content string
}
type activity struct{
  FileName string
  StudentName string
  FileImg string
  AType string
}
type list struct{
  Activities []*activity
  Atype string
}
func getpdf(test string) []*activity{
  files,err := ioutil.ReadDir("activities/"+test)
  if(err!=nil){
    log.Fatal(err)
  }
  var data []*activity
  for _, file := range files {
    var element activity
    element.AType= test
    element.FileName = file.Name()
    element.StudentName= strings.Replace(file.Name(),".pdf","",-1)
    element.FileImg= strings.Replace(file.Name(),".pdf",".jpg",-1)
    data = append(data,&element)
  }
  return data
}
func home(w http.ResponseWriter, r *http.Request){
  t,err := template.ParseFiles("index.html")

  if err!=nil{
    fmt.Fprintf(w,"error")
  }
  t.Execute(w,"test")
}
func blog(w http.ResponseWriter, r *http.Request){

  clientOptions :=options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
  client,err :=mongo.Connect(context.TODO(),clientOptions)
  if err!=nil{
    log.Fatal(err)
  }

  err = client.Ping(context.TODO(),nil)
  if err!=nil{
    log.Fatal(err)
  }
  collection := client.Database("cse").Collection("blogs")
  fmt.Println("connected to mongodb")

  var blogs []*Blog
  findOptions := options.Find()
  findOptions.SetSort(bson.D{{"$natural",-1}})
  cur,err := collection.Find(context.TODO(),bson.D{{}},findOptions)
  if err!=nil{
    log.Fatal(err)
  }

  for cur.Next(context.TODO()){

    var element Blog

    err := cur.Decode(&element)

    if err!=nil{
      log.Fatal(err)
    }

    blogs = append(blogs,&element)
  }

    if err:=cur.Err(); err!=nil{
        log.Fatal(err)
      }

    cur.Close(context.TODO())
    fmt.Println("multiple documents fetched")

  t,err := template.ParseFiles("blog.html")

  if err!=nil{
    fmt.Fprintf(w,"error")
  }
  type lists struct{
    Todo []*Blog
  }
  var data lists
  data.Todo=blogs
  t.Execute(w,data)
}
func achievement(w http.ResponseWriter, r *http.Request){
  t,err := template.ParseFiles("achievers.html")

  if err!=nil{
    fmt.Fprintf(w,"error")
  }
  t.Execute(w,nil)
}
func team(w http.ResponseWriter, r *http.Request){
  t,err := template.ParseFiles("team.html")

  if err!=nil{
    fmt.Fprintf(w,"error")
  }
  t.Execute(w,nil)
}
func events(w http.ResponseWriter, r *http.Request){
  t,err := template.ParseFiles("events.html")

  if err!=nil{
    fmt.Fprintf(w,"error")
  }
  t.Execute(w,nil)
}
func project(w http.ResponseWriter, r *http.Request){

  t,err := template.ParseFiles("project.html")
  if err!=nil{
    fmt.Fprintf(w,"error")
  }
  var Project list
  Project.Atype="Projects"
  Project.Activities = getpdf("projects")
  t.Execute(w,Project)
}
func seminars(w http.ResponseWriter, r *http.Request){

  t,err := template.ParseFiles("project.html")
  if err!=nil{
    fmt.Fprintf(w,"error")
  }
  var Seminar list
  Seminar.Atype = "Seminars"
  Seminar.Activities = getpdf("seminars")
  t.Execute(w,Seminar)
}

func writeblog(w http.ResponseWriter, r *http.Request){
    fmt.Println("called by get")
  t,err :=template.ParseFiles("writeablog.html")

  if err!=nil{
    fmt.Fprintf(w,"error")
  }
  t.Execute(w,nil)
}
func newblog(w http.ResponseWriter, r *http.Request){
      s := struct{
        Uploaded bool
      }{
        Uploaded:false,
      }
      //"mongodb://localhost:27017"
    clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
    client,err := mongo.Connect(context.TODO(),clientOptions)

    if err!=nil{
      log.Fatal(err)
    }

    err= client.Ping(context.TODO(),nil)

    if err!=nil{
      log.Fatal(err)
    }

    collection := client.Database("cse").Collection("blogs")
    fmt.Println("connected to mongodb")

    r.ParseForm()
    name:=r.FormValue("name")
    sub :=r.FormValue("subject")
    sem :=r.FormValue("semester")
    content:=r.FormValue("content")
    newBlog :=Blog{Name:name, Subject:sub, Sem:sem, Content:content}

    insert,err := collection.InsertOne(context.TODO(),newBlog)

    if err!=nil{
      log.Fatal(err)
    }else{
      s = struct{
            Uploaded bool
          }{
            Uploaded:true,
          }
      }

    fmt.Println("inserted : %v",insert.InsertedID)
    t,err :=template.ParseFiles("writeablog.html")

    if err!=nil{
      fmt.Fprintf(w,"error")
    }
    t.Execute(w,s)
}
func GetPort() string{
var port = os.Getenv("PORT")
 	// Set a default port if there is nothing in the environment
 	if port == "" {
		port = "8080"
    fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
 	}
  return ":" + port
}
func requestHandler(){
  m := mux.NewRouter().StrictSlash(true)
  m.PathPrefix("/css/").Handler(http.StripPrefix("/css/",http.FileServer(http.Dir("css"))))
  m.PathPrefix("/images/").Handler(http.StripPrefix("/images/",http.FileServer(http.Dir("images"))))
  m.PathPrefix("/js/").Handler(http.StripPrefix("/js/",http.FileServer(http.Dir("js"))))
  m.PathPrefix("/activities/").Handler(http.StripPrefix("/activities/",http.FileServer(http.Dir("activities"))))
  m.HandleFunc("/",home)
  m.HandleFunc("/blog",blog)
  m.HandleFunc("/achievers",achievement)
  m.HandleFunc("/events",events)
  m.HandleFunc("/project",project)
  m.HandleFunc("/seminar",seminars)
  m.HandleFunc("/team",team)
  m.HandleFunc("/writeablog",writeblog).Methods("GET")
  m.HandleFunc("/writeablog",newblog).Methods("POST")
  log.Fatal(http.ListenAndServe(GetPort(), m))
}
func main(){
  requestHandler()
}
