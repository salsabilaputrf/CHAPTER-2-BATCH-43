package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)


var Data = map[string]interface{}{
	"Title":   "Personal Web",
	"IsLogin": true,
}


type dataProject struct {
	Id int
	ProjectName	string
	StartDate string
	EndDate	string
	Description	string
	Technologies []string
	Duration string
}

var Projects = []dataProject{
	/*
	{
		ProjectName:  "Dumbways Mobile Apps",
		StartDate: "2022-12-02",
		EndDate:    "2022-12-29",
		Description:   "Aplikasi mobile dari dumbways",
		Technologies:   []string{"nodeJs", "reactJs", "nextJs", "typeScript"},
	},*/
}

// main function
func main() {
	route := mux.NewRouter()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.HandleFunc("/", helloWorld).Methods("GET")
	route.HandleFunc("/home", home).Methods("GET").Name("home")
	route.HandleFunc("/myProject", project).Methods("GET")
	route.HandleFunc("/myProject", addProject).Methods("POST")
	route.HandleFunc("/projectDetail/{id}", projectDetail).Methods("GET")
	route.HandleFunc("/deleteProject/{id}", deleteProject).Methods("GET")
	route.HandleFunc("/myProject", formEdit).Methods("GET")
	route.HandleFunc("/contact", contactMe).Methods("GET")


	// port := 5000
	fmt.Println("Server is running on port 5000")
	http.ListenAndServe("localhost:5000", route)
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello world!"))
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("view/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func project(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("view/my-project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	respData := map[string]interface{} {
		"Data":  Data,
        "Projects": Projects,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, respData)
}

func projectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("view/detail-project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	ProjectDetail := dataProject{}

	for index, data := range Projects{
		if index == id {
				ProjectDetail = dataProject{
					Id: id,
					ProjectName: data.ProjectName,
					StartDate: data.StartDate,
					EndDate: data.EndDate,
					Description: data.Description,
					Technologies: data.Technologies,
					Duration: data.Duration,
			}
		}
	}


	respDataDetail := map[string]interface{}{
		"Data": Data,
		"ProjectDetail": ProjectDetail,
	}

	fmt.Println("Id			  : ", ProjectDetail.Id)
	fmt.Println("Project Name : ", ProjectDetail.ProjectName)
	fmt.Println("Start Date   : ", ProjectDetail.StartDate)
	fmt.Println("End Date     : ", ProjectDetail.EndDate)
	fmt.Println("Description  : ", ProjectDetail.Description)
	fmt.Println("Technologies : ", ProjectDetail.Technologies)
	fmt.Println("Duration     : ", ProjectDetail.Duration)

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, respDataDetail)
}

func addProject (w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	
	projectName := r.PostForm.Get("projectName")
	startDate := r.PostForm.Get("startDate")
	endDate := r.PostForm.Get("endDate")
	desc := r.PostForm.Get("desc")
	tech := r.Form["technologi"]

	// Menghitung durasi
	// Parsing string to time.Time

	// Start Date
	startDateTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	// End Date
	endDateTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	// Perbedaan nya berupa : jam menit detik
	distance := endDateTime.Sub(startDateTime)

	// Menghitung durasi
	var duration string
	year := int(distance.Hours()/(12 * 30 * 24))
	if year != 0 {
		duration = strconv.Itoa(year) + " tahun"
	}else{
		month := int(distance.Hours()/(30 * 24))
		if month != 0 {
			duration = strconv.Itoa(month) + " bulan"
		}else{
			week := int(distance.Hours()/(7 * 24))
			if week != 0 {
				duration = strconv.Itoa(week) + " minggu"
			}else{
				day := int(distance.Hours()/(24))
				if day != 0 {
					duration = strconv.Itoa(day) + " hari"
				}
			}
		}
	}

	var newProject = dataProject {
		ProjectName: projectName,
		StartDate: startDate,
		EndDate: endDate,
		Description: desc,
		Technologies: tech,
		Duration: duration,
	}
	/*  -- Untuk menampilkan di console ( Task Day 7 ) --
	fmt.Println("Project Name : ", newProject.ProjectName)
	fmt.Println("Start Date   : ", newProject.StartDate)
	fmt.Println("End Date     : ", newProject.EndDate)
	fmt.Println("Description  : ", newProject.Description)
	fmt.Println("Technologies : ", newProject.Technologies)
	fmt.Println("Duration     : ", newProject.Duration)
	*/

	Projects = append(Projects, newProject)

	http.Redirect(w, r, "/myProject", http.StatusMovedPermanently)
}

func deleteProject (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	Projects = append(Projects[:id], Projects[id+1:]...)

	http.Redirect(w, r, "/myProject", http.StatusMovedPermanently)

}

func formEdit (w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("view/my-project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	respData := map[string]interface{} {
		"Data":  Data,
        "Projects": Projects,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, respData)
}

func contactMe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("view/contact-form.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

