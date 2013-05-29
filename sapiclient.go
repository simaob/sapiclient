package main

import (
  "fmt"
  "encoding/json"
  "github.com/bmizerany/pat"
  "net/http"
  _ "github.com/lxn/go-pgsql" //or _ "github.com/bmizerany/pq"
  "database/sql"
  "log"
)

type taxonConcept struct {
  Id string
  Full_name string
  Kingdom string
}


func taxonConceptsHandler (w http.ResponseWriter, r *http.Request) {
  the_id := r.URL.Query().Get(":id")

  db, err := sql.Open("postgres", "user=postgres dbname=sapi_development password=postgres, sslmode=disable")

  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

  var taxon_concept taxonConcept

  if the_id != "all" {
    err = db.QueryRow("SELECT id, full_name, kingdom_name FROM taxon_concepts_mview WHERE id=$1;", the_id).Scan(&taxon_concept.Id, &taxon_concept.Full_name, &taxon_concept.Kingdom)

    if err != nil {
      log.Fatal(err)
    }

    taxon_json, err := json.Marshal(taxon_concept)
    if err == nil {
    w.Write(taxon_json)
    fmt.Println(taxon_concept)
    } else {
      log.Fatal(err)
    }
  } else {
    rows, _ := db.Query("SELECT id, full_name, kingdom_name FROM taxon_concepts_mview LIMIT 15;")
    var taxon_concepts []taxonConcept
    for rows.Next() {
      if err := rows.Scan(&taxon_concept.Id, &taxon_concept.Full_name, &taxon_concept.Kingdom); err != nil {
            log.Fatal(err)
      }
      taxon_concepts = append(taxon_concepts, taxon_concept)
    }
    taxon_json, err := json.Marshal(taxon_concepts)
    if err == nil {
      w.Write(taxon_json)
    } else {
      log.Fatal(err)
    }
  }
}

func main () {

  matcher := pat.New()

  matcher.Get("/taxon_concepts/:id", http.HandlerFunc(taxonConceptsHandler))

  http.Handle("/taxon_concepts/", matcher)
  http.ListenAndServe(":4000", nil)
}
