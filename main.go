package main 
import	(
	"net/http"
	"log")

func main(){
	sermux := http.NewServeMux()
	s:=&http.Server{
		Addr:           ":8080",
		Handler:        sermux,
	}
	sermux.Handle("/",http.FileServer(http.Dir("./")))

	err := s.ListenAndServe()
	if err != nil{
		log.Printf("Error occured")
		return
	}
	return
}
