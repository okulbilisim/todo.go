package main
import (
	"fmt"
	"flag"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)
func main(){
	list := flag.Bool("list",false,"List tasks")
	newTask := flag.String("new","","New Task")
	deleteTask := flag.Int("delete",0,"Delete Task")
	flag.Parse()

	
	db, err := sql.Open("sqlite3","todo.db")
	if err!=nil {
		panic(err)
	}
	defer db.Close()
	table(db)

	switch{
	case *list:
		listAll(db)
	case *newTask!="":
		createTask(newTask, db)
	case *deleteTask!=0:
		removeTask(deleteTask, db)
	default:
		help()
	}

}
func removeTask(id *int, db *sql.DB){
	tx, err:=db.Begin()
	if err!=nil {
		panic(err)
	}
	stmt, err := tx.Prepare("DELETE FROM todogo WHERE id=?")
	defer stmt.Close()
	if err!=nil{
		panic(err)
	}
	stmt.Exec(id)
	tx.Commit()
	fmt.Printf("Deleted #%d",*id)
}

func listAll(db *sql.DB){
	rows, err := db.Query("select id,title from todogo")
	if err!=nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var title string
		rows.Scan(&id,&title)
		fmt.Printf("#%d %s \n",id,title)
	}
}

func createTask(newTask *string, db *sql.DB){
	tx, err := db.Begin()
	if err!=nil {
		panic(err)
	}
	stmt, err := tx.Prepare("Insert Into todogo(id,title) values(?,?)")
	defer stmt.Close()
	if err!=nil {
		panic(err)
	}
	res, err := stmt.Exec(nil, *newTask)
	tx.Commit()
	if err!=nil {
		panic(err)
	}
	id,err:=res.LastInsertId();
	if err!=nil{
		panic(err)
	}
	fmt.Println("Task created. ",id)
}
func help(){
	fmt.Println(`
### TODO.go 

--help 		Help screen
--new=[task]  	Create new task
--list  	List all tasks
--delete=[id] 	Delete task
		`)
}
func table(db *sql.DB) {
	var query string;
	query =`
		CREATE TABLE IF NOT EXISTS todogo
		(id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		 title TEXT NOT NULL);`

	_,err := db.Exec(query)
	if err!=nil {
		panic(err)
	}
}