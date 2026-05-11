package main

import (
	"database/sql"
	"fmt"
	"sync"
	"time"
)

var length = 840

func main() {

	var wg *sync.WaitGroup = new(sync.WaitGroup)
	wg.Add(length)
	usersNmes := []string{"Elena Tapia", "Samir Austin", "Alivia Graham", "Giovanni Ingram",
		"Katie Vaughan", "Castiel Warner", "Wynter Esparza", "Carl Waters", "Bristol Richards",
		"Holden Lamb", "Amaia Moore", "Levi Contreras", "Daniela Peters", "Patrick Cobb",
		"Aviana Chapman", "Knox Cuevas", "Adele Atkinson", "Duke Matthews", "Lila Sharp",
		"Royce Blackwell", "Saoirse Lopez", "Michael Barajas", "Keilani House", "Yehuda Shannon",
		"Harlee Fletcher", "Jay Griffin", "Charlie Norman", "Aziel Farmer", "Madelynn Barry",
		"Emery Atkinson", "Jazmin Lam", "Bodie Jaramillo", "Guadalupe Jackson", "Sebastian Casey",
		"Sylvia Chapman", "Knox Hardy", "Jessica Grimes", "Harlan McBride", "Kelsey Jones",
		"William Villegas", "Jessie Hail", "Hector Wilkins", "Amalia Olson", "Malachi Chang",
		"Ophelia Orozco", "Keanu Gonzales", "Hadley Nava", "Stefan Boyd", "Georgia Welch",
		"Hendrix Lozano", "Cecelia Moody", "Ryland Nielsen", "Vienna Hodges", "Alonzo Cisneros",
		"Janelle Ruiz", "Austin Gutierrez", "Savannah Hamilton", "Jason Brock", "Jada Rodriguez",
		"Henry Collier", "Ivory West", "Diego Skinner", "Mara Huerta", "Douglas Weaver", "Teagan Clark", "John Morton", "Mallory Vega", "Aidan Bean", "Jenesis Garner", "Sage Richardson", "Allison Powell", "Bennett Moon", "Naya Fitzpatrick", "Blaze Glenn", "Blaire White", "Aiden Daniel", "Joy Wilson", "Daniel Porter", "Ryleigh McKee", "Bjorn Sexton", "Ellen Chang", "Ari Kelly", "Ruby Wagner", "Enzo Hogan", "Kathryn Rowe", "Kamden Beck", "Gia Andrews", "Lukas Ventura", "Zora Larson", "Rafael Reeves", "Lana Butler", "Ryder Miller", "Isabella Cantrell", "Harris Guevara", "Teresa Conley", "Marvin Stout", "Chana Matthews", "Preston Benton", "Anais Peters", "Patrick Reese", "Rosemary Conley", "Marvin Dunlap", "Iliana Patel", "Parker Donaldson", "Natasha McConnell", "London Kemp", "Anika Thornton", "Malik Glass", "Clare Weiss", "Koa Russo", "Tinsley Anderson", "Jacob Keller", "Logan Griffith", "Franklin Rojas", "Adaline Shaw", "Elliot Briggs", "Alia McConnell", "London Franklin", "Angela Shepherd", "Ronald York",
		"Elena Tapia", "Samir Austin", "Alivia Graham", "Giovanni Ingram",
		"Katie Vaughan", "Castiel Warner", "Wynter Esparza", "Carl Waters", "Bristol Richards",
		"Holden Lamb", "Amaia Moore", "Levi Contreras", "Daniela Peters", "Patrick Cobb",
		"Aviana Chapman", "Knox Cuevas", "Adele Atkinson", "Duke Matthews", "Lila Sharp",
		"Royce Blackwell", "Saoirse Lopez", "Michael Barajas", "Keilani House", "Yehuda Shannon",
		"Harlee Fletcher", "Jay Griffin", "Charlie Norman", "Aziel Farmer", "Madelynn Barry",
		"Emery Atkinson", "Jazmin Lam", "Bodie Jaramillo", "Guadalupe Jackson", "Sebastian Casey",
		"Sylvia Chapman", "Knox Hardy", "Jessica Grimes", "Harlan McBride", "Kelsey Jones",
		"William Villegas", "Jessie Hail", "Hector Wilkins", "Amalia Olson", "Malachi Chang",
		"Ophelia Orozco", "Keanu Gonzales", "Hadley Nava", "Stefan Boyd", "Georgia Welch",
		"Hendrix Lozano", "Cecelia Moody", "Ryland Nielsen", "Vienna Hodges", "Alonzo Cisneros",
		"Janelle Ruiz", "Austin Gutierrez", "Savannah Hamilton", "Jason Brock", "Jada Rodriguez",
		"Henry Collier", "Ivory West", "Diego Skinner", "Mara Huerta", "Douglas Weaver", "Teagan Clark", "John Morton", "Mallory Vega", "Aidan Bean", "Jenesis Garner", "Sage Richardson", "Allison Powell", "Bennett Moon", "Naya Fitzpatrick", "Blaze Glenn", "Blaire White", "Aiden Daniel", "Joy Wilson", "Daniel Porter", "Ryleigh McKee", "Bjorn Sexton", "Ellen Chang", "Ari Kelly", "Ruby Wagner", "Enzo Hogan", "Kathryn Rowe", "Kamden Beck", "Gia Andrews", "Lukas Ventura", "Zora Larson", "Rafael Reeves", "Lana Butler", "Ryder Miller", "Isabella Cantrell", "Harris Guevara", "Teresa Conley", "Marvin Stout", "Chana Matthews", "Preston Benton", "Anais Peters", "Patrick Reese", "Rosemary Conley", "Marvin Dunlap", "Iliana Patel", "Parker Donaldson", "Natasha McConnell", "London Kemp", "Anika Thornton", "Malik Glass", "Clare Weiss", "Koa Russo", "Tinsley Anderson", "Jacob Keller", "Logan Griffith", "Franklin Rojas", "Adaline Shaw", "Elliot Briggs", "Alia McConnell", "London Franklin", "Angela Shepherd", "Ronald York",
		"Elena Tapia", "Samir Austin", "Alivia Graham", "Giovanni Ingram",
		"Katie Vaughan", "Castiel Warner", "Wynter Esparza", "Carl Waters", "Bristol Richards",
		"Holden Lamb", "Amaia Moore", "Levi Contreras", "Daniela Peters", "Patrick Cobb",
		"Aviana Chapman", "Knox Cuevas", "Adele Atkinson", "Duke Matthews", "Lila Sharp",
		"Royce Blackwell", "Saoirse Lopez", "Michael Barajas", "Keilani House", "Yehuda Shannon",
		"Harlee Fletcher", "Jay Griffin", "Charlie Norman", "Aziel Farmer", "Madelynn Barry",
		"Emery Atkinson", "Jazmin Lam", "Bodie Jaramillo", "Guadalupe Jackson", "Sebastian Casey",
		"Sylvia Chapman", "Knox Hardy", "Jessica Grimes", "Harlan McBride", "Kelsey Jones",
		"William Villegas", "Jessie Hail", "Hector Wilkins", "Amalia Olson", "Malachi Chang",
		"Ophelia Orozco", "Keanu Gonzales", "Hadley Nava", "Stefan Boyd", "Georgia Welch",
		"Hendrix Lozano", "Cecelia Moody", "Ryland Nielsen", "Vienna Hodges", "Alonzo Cisneros",
		"Janelle Ruiz", "Austin Gutierrez", "Savannah Hamilton", "Jason Brock", "Jada Rodriguez",
		"Henry Collier", "Ivory West", "Diego Skinner", "Mara Huerta", "Douglas Weaver", "Teagan Clark", "John Morton", "Mallory Vega", "Aidan Bean", "Jenesis Garner", "Sage Richardson", "Allison Powell", "Bennett Moon", "Naya Fitzpatrick", "Blaze Glenn", "Blaire White", "Aiden Daniel", "Joy Wilson", "Daniel Porter", "Ryleigh McKee", "Bjorn Sexton", "Ellen Chang", "Ari Kelly", "Ruby Wagner", "Enzo Hogan", "Kathryn Rowe", "Kamden Beck", "Gia Andrews", "Lukas Ventura", "Zora Larson", "Rafael Reeves", "Lana Butler", "Ryder Miller", "Isabella Cantrell", "Harris Guevara", "Teresa Conley", "Marvin Stout", "Chana Matthews", "Preston Benton", "Anais Peters", "Patrick Reese", "Rosemary Conley", "Marvin Dunlap", "Iliana Patel", "Parker Donaldson", "Natasha McConnell", "London Kemp", "Anika Thornton", "Malik Glass", "Clare Weiss", "Koa Russo", "Tinsley Anderson", "Jacob Keller", "Logan Griffith", "Franklin Rojas", "Adaline Shaw", "Elliot Briggs", "Alia McConnell", "London Franklin", "Angela Shepherd", "Ronald York",
		"Elena Tapia", "Samir Austin", "Alivia Graham", "Giovanni Ingram",
		"Katie Vaughan", "Castiel Warner", "Wynter Esparza", "Carl Waters", "Bristol Richards",
		"Holden Lamb", "Amaia Moore", "Levi Contreras", "Daniela Peters", "Patrick Cobb",
		"Aviana Chapman", "Knox Cuevas", "Adele Atkinson", "Duke Matthews", "Lila Sharp",
		"Royce Blackwell", "Saoirse Lopez", "Michael Barajas", "Keilani House", "Yehuda Shannon",
		"Harlee Fletcher", "Jay Griffin", "Charlie Norman", "Aziel Farmer", "Madelynn Barry",
		"Emery Atkinson", "Jazmin Lam", "Bodie Jaramillo", "Guadalupe Jackson", "Sebastian Casey",
		"Sylvia Chapman", "Knox Hardy", "Jessica Grimes", "Harlan McBride", "Kelsey Jones",
		"William Villegas", "Jessie Hail", "Hector Wilkins", "Amalia Olson", "Malachi Chang",
		"Ophelia Orozco", "Keanu Gonzales", "Hadley Nava", "Stefan Boyd", "Georgia Welch",
		"Hendrix Lozano", "Cecelia Moody", "Ryland Nielsen", "Vienna Hodges", "Alonzo Cisneros",
		"Janelle Ruiz", "Austin Gutierrez", "Savannah Hamilton", "Jason Brock", "Jada Rodriguez",
		"Henry Collier", "Ivory West", "Diego Skinner", "Mara Huerta", "Douglas Weaver", "Teagan Clark", "John Morton", "Mallory Vega", "Aidan Bean", "Jenesis Garner", "Sage Richardson", "Allison Powell", "Bennett Moon", "Naya Fitzpatrick", "Blaze Glenn", "Blaire White", "Aiden Daniel", "Joy Wilson", "Daniel Porter", "Ryleigh McKee", "Bjorn Sexton", "Ellen Chang", "Ari Kelly", "Ruby Wagner", "Enzo Hogan", "Kathryn Rowe", "Kamden Beck", "Gia Andrews", "Lukas Ventura", "Zora Larson", "Rafael Reeves", "Lana Butler", "Ryder Miller", "Isabella Cantrell", "Harris Guevara", "Teresa Conley", "Marvin Stout", "Chana Matthews", "Preston Benton", "Anais Peters", "Patrick Reese", "Rosemary Conley", "Marvin Dunlap", "Iliana Patel", "Parker Donaldson", "Natasha McConnell", "London Kemp", "Anika Thornton", "Malik Glass", "Clare Weiss", "Koa Russo", "Tinsley Anderson", "Jacob Keller", "Logan Griffith", "Franklin Rojas", "Adaline Shaw", "Elliot Briggs", "Alia McConnell", "London Franklin", "Angela Shepherd", "Ronald York",
		"Elena Tapia", "Samir Austin", "Alivia Graham", "Giovanni Ingram",
		"Katie Vaughan", "Castiel Warner", "Wynter Esparza", "Carl Waters", "Bristol Richards",
		"Holden Lamb", "Amaia Moore", "Levi Contreras", "Daniela Peters", "Patrick Cobb",
		"Aviana Chapman", "Knox Cuevas", "Adele Atkinson", "Duke Matthews", "Lila Sharp",
		"Royce Blackwell", "Saoirse Lopez", "Michael Barajas", "Keilani House", "Yehuda Shannon",
		"Harlee Fletcher", "Jay Griffin", "Charlie Norman", "Aziel Farmer", "Madelynn Barry",
		"Emery Atkinson", "Jazmin Lam", "Bodie Jaramillo", "Guadalupe Jackson", "Sebastian Casey",
		"Sylvia Chapman", "Knox Hardy", "Jessica Grimes", "Harlan McBride", "Kelsey Jones",
		"William Villegas", "Jessie Hail", "Hector Wilkins", "Amalia Olson", "Malachi Chang",
		"Ophelia Orozco", "Keanu Gonzales", "Hadley Nava", "Stefan Boyd", "Georgia Welch",
		"Hendrix Lozano", "Cecelia Moody", "Ryland Nielsen", "Vienna Hodges", "Alonzo Cisneros",
		"Janelle Ruiz", "Austin Gutierrez", "Savannah Hamilton", "Jason Brock", "Jada Rodriguez",
		"Henry Collier", "Ivory West", "Diego Skinner", "Mara Huerta", "Douglas Weaver", "Teagan Clark", "John Morton", "Mallory Vega", "Aidan Bean", "Jenesis Garner", "Sage Richardson", "Allison Powell", "Bennett Moon", "Naya Fitzpatrick", "Blaze Glenn", "Blaire White", "Aiden Daniel", "Joy Wilson", "Daniel Porter", "Ryleigh McKee", "Bjorn Sexton", "Ellen Chang", "Ari Kelly", "Ruby Wagner", "Enzo Hogan", "Kathryn Rowe", "Kamden Beck", "Gia Andrews", "Lukas Ventura", "Zora Larson", "Rafael Reeves", "Lana Butler", "Ryder Miller", "Isabella Cantrell", "Harris Guevara", "Teresa Conley", "Marvin Stout", "Chana Matthews", "Preston Benton", "Anais Peters", "Patrick Reese", "Rosemary Conley", "Marvin Dunlap", "Iliana Patel", "Parker Donaldson", "Natasha McConnell", "London Kemp", "Anika Thornton", "Malik Glass", "Clare Weiss", "Koa Russo", "Tinsley Anderson", "Jacob Keller", "Logan Griffith", "Franklin Rojas", "Adaline Shaw", "Elliot Briggs", "Alia McConnell", "London Franklin", "Angela Shepherd", "Ronald York",
		"Elena Tapia", "Samir Austin", "Alivia Graham", "Giovanni Ingram",
		"Katie Vaughan", "Castiel Warner", "Wynter Esparza", "Carl Waters", "Bristol Richards",
		"Holden Lamb", "Amaia Moore", "Levi Contreras", "Daniela Peters", "Patrick Cobb",
		"Aviana Chapman", "Knox Cuevas", "Adele Atkinson", "Duke Matthews", "Lila Sharp",
		"Royce Blackwell", "Saoirse Lopez", "Michael Barajas", "Keilani House", "Yehuda Shannon",
		"Harlee Fletcher", "Jay Griffin", "Charlie Norman", "Aziel Farmer", "Madelynn Barry",
		"Emery Atkinson", "Jazmin Lam", "Bodie Jaramillo", "Guadalupe Jackson", "Sebastian Casey",
		"Sylvia Chapman", "Knox Hardy", "Jessica Grimes", "Harlan McBride", "Kelsey Jones",
		"William Villegas", "Jessie Hail", "Hector Wilkins", "Amalia Olson", "Malachi Chang",
		"Ophelia Orozco", "Keanu Gonzales", "Hadley Nava", "Stefan Boyd", "Georgia Welch",
		"Hendrix Lozano", "Cecelia Moody", "Ryland Nielsen", "Vienna Hodges", "Alonzo Cisneros",
		"Janelle Ruiz", "Austin Gutierrez", "Savannah Hamilton", "Jason Brock", "Jada Rodriguez",
		"Henry Collier", "Ivory West", "Diego Skinner", "Mara Huerta", "Douglas Weaver", "Teagan Clark", "John Morton", "Mallory Vega", "Aidan Bean", "Jenesis Garner", "Sage Richardson", "Allison Powell", "Bennett Moon", "Naya Fitzpatrick", "Blaze Glenn", "Blaire White", "Aiden Daniel", "Joy Wilson", "Daniel Porter", "Ryleigh McKee", "Bjorn Sexton", "Ellen Chang", "Ari Kelly", "Ruby Wagner", "Enzo Hogan", "Kathryn Rowe", "Kamden Beck", "Gia Andrews", "Lukas Ventura", "Zora Larson", "Rafael Reeves", "Lana Butler", "Ryder Miller", "Isabella Cantrell", "Harris Guevara", "Teresa Conley", "Marvin Stout", "Chana Matthews", "Preston Benton", "Anais Peters", "Patrick Reese", "Rosemary Conley", "Marvin Dunlap", "Iliana Patel", "Parker Donaldson", "Natasha McConnell", "London Kemp", "Anika Thornton", "Malik Glass", "Clare Weiss", "Koa Russo", "Tinsley Anderson", "Jacob Keller", "Logan Griffith", "Franklin Rojas", "Adaline Shaw", "Elliot Briggs", "Alia McConnell", "London Franklin", "Angela Shepherd", "Ronald York",
		"Elena Tapia", "Samir Austin", "Alivia Graham", "Giovanni Ingram",
		"Katie Vaughan", "Castiel Warner", "Wynter Esparza", "Carl Waters", "Bristol Richards",
		"Holden Lamb", "Amaia Moore", "Levi Contreras", "Daniela Peters", "Patrick Cobb",
		"Aviana Chapman", "Knox Cuevas", "Adele Atkinson", "Duke Matthews", "Lila Sharp",
		"Royce Blackwell", "Saoirse Lopez", "Michael Barajas", "Keilani House", "Yehuda Shannon",
		"Harlee Fletcher", "Jay Griffin", "Charlie Norman", "Aziel Farmer", "Madelynn Barry",
		"Emery Atkinson", "Jazmin Lam", "Bodie Jaramillo", "Guadalupe Jackson", "Sebastian Casey",
		"Sylvia Chapman", "Knox Hardy", "Jessica Grimes", "Harlan McBride", "Kelsey Jones",
		"William Villegas", "Jessie Hail", "Hector Wilkins", "Amalia Olson", "Malachi Chang",
		"Ophelia Orozco", "Keanu Gonzales", "Hadley Nava", "Stefan Boyd", "Georgia Welch",
		"Hendrix Lozano", "Cecelia Moody", "Ryland Nielsen", "Vienna Hodges", "Alonzo Cisneros",
		"Janelle Ruiz", "Austin Gutierrez", "Savannah Hamilton", "Jason Brock", "Jada Rodriguez",
		"Henry Collier", "Ivory West", "Diego Skinner", "Mara Huerta", "Douglas Weaver", "Teagan Clark", "John Morton", "Mallory Vega", "Aidan Bean", "Jenesis Garner", "Sage Richardson", "Allison Powell", "Bennett Moon", "Naya Fitzpatrick", "Blaze Glenn", "Blaire White", "Aiden Daniel", "Joy Wilson", "Daniel Porter", "Ryleigh McKee", "Bjorn Sexton", "Ellen Chang", "Ari Kelly", "Ruby Wagner", "Enzo Hogan", "Kathryn Rowe", "Kamden Beck", "Gia Andrews", "Lukas Ventura", "Zora Larson", "Rafael Reeves", "Lana Butler", "Ryder Miller", "Isabella Cantrell", "Harris Guevara", "Teresa Conley", "Marvin Stout", "Chana Matthews", "Preston Benton", "Anais Peters", "Patrick Reese", "Rosemary Conley", "Marvin Dunlap", "Iliana Patel", "Parker Donaldson", "Natasha McConnell", "London Kemp", "Anika Thornton", "Malik Glass", "Clare Weiss", "Koa Russo", "Tinsley Anderson", "Jacob Keller", "Logan Griffith", "Franklin Rojas", "Adaline Shaw", "Elliot Briggs", "Alia McConnell", "London Franklin", "Angela Shepherd", "Ronald York",
	}

	db := InitializeDB(usersNmes)
	if db == nil {
		fmt.Println("Error initializing DB")
		return
	}
	handler := flightSeatHandler{
		flightDB: db,
	}

	pnrs := make([]string, 0, length)
	x, err := db.Query("SELECT PNR FROM PLANE_USERS")
	if err != nil {
		fmt.Println(err)
		return
	}

	for x.Next() {
		var pnr string
		x.Scan(&pnr)
		pnrs = append(pnrs, pnr)
	}

	t1 := time.Now()

	for i := 0; i < length; i++ {
		go func() {
			handler.AllocateSeats(pnrs[i], wg)
		}()
	}

	wg.Wait()
	fmt.Println(time.Since(t1))
	printData(db)
	ExecuteShutdownMethods(db)

}

func ExecuteShutdownMethods(db *sql.DB) {
	db.Exec("DROP TABLE IF EXISTS PLANE;")
	db.Exec("DROP TABLE IF EXISTS PLANE_USERS;")
}

func printData(db *sql.DB) {
	rows, _ := db.Query("SELECT * FROM PLANE")
	count := 0
	var setID, userID string
	for rows.Next() {
		rows.Scan(&setID, &userID)
		fmt.Println(setID, "--------------", userID)
		count++
	}
	fmt.Println(count)
}
