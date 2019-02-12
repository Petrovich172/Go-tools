package main

import (
	"log"
	"github.com/go-pg/pg"
	"math/rand"
	// "encoding/json"
	// "io/ioutil"
	"io"
	"os"
	"encoding/xml"
	utils "./pkgs/utils"
	"./pkgs/cfg"

	// inits "pjob/pkgs/init"

	// "github.com/gin-gonic/gin"
	// "github.com/go-redis/redis"
)


// our struct which contains the complete
// array of all Users in the file
type Edge struct {
	// XMLName xml.Name `xml:"users"`
	Id	string	`xml:"id"	sql:"id"`
	Oneway	string	`xml:"oneway"	sql:"oneway"`
	Surface	string	`xml:"surface"	sql:"surface"`
	Highway	string	`xml:"highway"	sql:"highway"`
	Geom	utils.MultiLineString	`xml:"geom"	sql:"geom"`
	// Users   []User   `xml:"user"`
}

// the user struct, this contains our
// Type attribute, our user's name and
// a social struct which will contain all
// our social links
type User struct {
	XMLName xml.Name `xml:"user"`
	Type    string   `xml:"type,attr"`
	Name    string   `xml:"name"`
	Social  Social   `xml:"social"`
}

// a simple struct which contains all our
// social links
type Social struct {
	XMLName  xml.Name `xml:"social"`
	Facebook string   `xml:"facebook"`
	Twitter  string   `xml:"twitter"`
	Youtube  string   `xml:"youtube"`
}


func main() {
	log.Println("Heey!")

	db := pg.Connect(&pg.Options{
			Addr:      "172.20.12.159" + ":" + "5432",
			User:      "postgres",
			Password:  "postgres",
			Database:  "postgres",
		})
	defer db.Close()

	// Open our xmlFile
	xmlFile, err := os.Open("sample.xml")
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Println(err)
	}

	log.Println("Successfully Opened sample.xml")
	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	// read our opened xmlFile as a byte array.
	
	// byteValue, _ := ioutil.ReadAll(xmlFile)

	// we initialize our Users array
	// var data []Edge
	var xmlData cfg.Map
	log.Println(xmlData)

	dbData := getSomeData(db)
	generate := rand.New(rand.NewSource(99)).Int63
	var nodeId cfg.Elem

	for i := 0; i < len(dbData); i++ {
		// log.Println("rand:", nodeId)
		// log.Println("smth:","\n",dbData[i].Geom)
		// log.Println("smth else:","\n",dbData[i].Geom.Coordinates[0])

		node := dbData[i].Geom.Coordinates[0]
		var nodeIDs []int64
		for y := 0; y < len(node); y++ {
			nodeId.ID = generate()
			xmlData.Nodes = append(	xmlData.Nodes, cfg.Node{
				Elem:	nodeId,
				Lat:	node[y][0],
				Lng:	node[y][1],
				}	)
			nodeIDs = append(nodeIDs, nodeId.ID)
		}

		for yy := 0; yy < len(nodeIDs); yy++ {
		// for yy := 0; yy < len(nodeIDs); yy++ {
			var tempId make([]struct {	ID int64 `xml:"ref,attr"`}, 3)
			tempId[yy].ID = nodeIDs[yy]
			log.Println(tempId)
			// log.Println("xmlData.Ways[i].Nds:",xmlData.Ways[i].Nds)
			// xmlData.Ways[i].Nds[0].ID = tempId.ID
			xmlData.Ways = append(xmlData.Ways, cfg.Way{
				Nds:	tempId,
			}	)
		}
			// xmlData.Ways[i].Nds = append( xmlData.Ways[i].Nds, struct{
				// ID:	nodeIDs[yy],},
				// }	)
			
	}

	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'users' which we defined above
	// xml.Unmarshal(	getSomeData(db), &data	)	

	// for i := 0; i < len(data); i++ {
	// 	log.Println("User Type: " + data[i].Id)
	// 	log.Println("User Name: " + data.Users[i].Highway)
	// 	log.Println("Facebook Url: " + data.Users[i].Geom)
	// }

	// creating output xml file
	f, err := os.Create("out.xml")
	if err != nil { panic(err) }
	defer f.Close()
	newFile := io.Writer(f)
	enc := xml.NewEncoder(newFile)
	enc.Indent("  ", "    ")
    	if err := enc.Encode(&xmlData); err != nil {
				log.Printf("error: %v\n", err)
		}
}

func getSomeData(db *pg.DB) []Edge {
	var ret []Edge
	var err error
	sqlString := "SELECT id, ST_AsGeoJSON(geom) as geom, oneway, surface, highway FROM public.tline_smaller limit 5"
	_, err = db.Model().Query(&ret, sqlString)
	if err != nil {
		log.Println("some shit happend:", "\n", err)		
	}
	log.Println("query:","\n",ret)
	// tmp := json.Marshal(ret)
	return ret
}