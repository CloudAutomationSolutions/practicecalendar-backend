package function

// User ...
// User definition to keep in the database
type User struct {
	ID       string    `json:"id" firestore:"id"`
	Projects []Project `json:"projects" firestore:"projects"`
}

// Project ...
// A colection of tasks and their completion dates and start date
type Project struct {
	ID             string   `json:"id" firestore:"id"`
	Name           string   `json:"name" firestore:"name"`
	Tasks          []Task   `json:"tasks" firestore:"tasks"`
	StartDate      string   `json:"startDate" firestore:"startDate"`
	CompletedDates []string `json:"completedDates" firestore:"completedDates"`
}

// Task ...
// Unit to be completed. Multiple tasks form a Project.
type Task struct {
	ID           string `json:"id" firestore:"id"`
	Name         string `json:"name" firestore:"name"`
	LastDoneDate string `json:"lastDoneDate" firestore:"lastDoneDate"`
}
