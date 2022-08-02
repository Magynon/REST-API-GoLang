package gateways

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"workshop/domain"

	"github.com/emicklei/go-restful/v3"
)

const elementsPatch = "/elements"
const storeURL = "http://localhost:8082/STORE/elements"

type API struct {
	users map[int]domain.Element
}

func NewAPI() *API {
	return &API{
		users: make(map[int]domain.Element),
	}
}

func (api *API) RegisterRoutes(ws *restful.WebService) {
	ws.Path("/API")
	// ws.Route(ws.POST(echoPath).To(api.echoPOSTHandler).Reads(restful.MIME_JSON).Writes(restful.MIME_JSON).Doc("Writes back a json with what you gave it"))
	// ws.Route(ws.GET(echoPath).To(api.echoGETHandler).Writes(restful.MIME_JSON).Doc("Writes back a json with what you gave it"))

	ws.Route(ws.GET(elementsPatch).To(api.GETHandler).Writes(restful.MIME_JSON).Doc("Writes back a json with what you gave it"))
	ws.Route(ws.POST(elementsPatch).To(api.POSTHandler).Writes(restful.MIME_JSON).Doc("Writes back a json with what you gave it"))
	// ws.Route(ws.PATCH(usersPatch).To(api.updateElement).Writes(restful.MIME_JSON).Doc("Writes back a json with what you gave it"))
	// ws.Route(ws.DELETE(usersPatch).To(api.deleteElement).Writes(restful.MIME_JSON).Doc("Writes back a json with what you gave it"))

}

func init() {
    rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    return string(b)
}

func (api *API) POSTHandler(req *restful.Request, resp *restful.Response) {
	body := req.Request.Body
	if body == nil {
		log.Printf("[ERROR] Couldn't read request body")
		resp.WriteServiceError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, "nil body"))
		return

	}
	defer body.Close()
	var err error
	if err != nil {
		log.Printf("[ERROR] Couldn't read request body")
		resp.WriteServiceError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()))
		return
	}
	data, err := ioutil.ReadAll(body)
	if err != nil {
		log.Printf("[ERROR] Couldn't read request body")
		resp.WriteServiceError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()))
		return
	}

	// declare array of products
	var elems = []domain.Element{}
	ok := json.Unmarshal(data, &elems)

	if ok != nil {
		log.Printf("[ERROR] Couldn't unmarshal body data")
	}

	client := &http.Client {}
	wg := &sync.WaitGroup{}

	for _, elem := range elems {
		fmt.Println(elem.Name)

		if elem.Name == "pepsi" {
			// nu adaug la db
			continue
		} else {
			wg.Add(1)

			elem.Id = RandStringRunes(10)
			go func(elem domain.Element) {
				defer wg.Done()
				marshaled_elem, err := json.Marshal(elem)

				fmt.Println(string(marshaled_elem))

				req, err := http.NewRequest("POST", storeURL, strings.NewReader(string(marshaled_elem)))

				if err != nil {
					fmt.Println(err)
					return
				  }
				  req.Header.Add("Content-Type", "application/json")
				
				  res, err := client.Do(req)
				  if err != nil {
					fmt.Println(err)
					return
				  }
				  defer res.Body.Close()

				  body, err := ioutil.ReadAll(res.Body)
				  if err != nil {
				    fmt.Println(err)
				    return
				  }
 				  fmt.Println(string(body))
			}(elem)
		}
	}

	wg.Wait()
	resp.WriteAsJson(elems)
}

func (api *API) GETHandler(req *restful.Request, resp *restful.Response) {
	ids := req.QueryParameters("id")
	if ids == nil {
		log.Printf("[ERROR] Failed to read id")
		resp.WriteError(http.StatusBadRequest, fmt.Errorf("element id must be provided"))
		return
	}

	client := &http.Client {}
	fmt.Println(ids)
	wg := &sync.WaitGroup{}
	var mu sync.Mutex
	var final_response string = ""

	for _, elem := range ids {
		wg.Add(1)

		go func(elem string) {
			defer wg.Done()
			requst, err := http.NewRequest("GET", storeURL, strings.NewReader(elem))
  
			if err != nil {
				fmt.Println(err)
				return
			}
			requst.Header.Add("Content-Type", "text/plain")
	
			res, err := client.Do(requst)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer res.Body.Close()
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(body))
			mu.Lock()
			final_response += string(body) + "\n"
			mu.Unlock()
		}(elem)
	}
	
	wg.Wait()
	resp.WriteAsJson(final_response)
}
