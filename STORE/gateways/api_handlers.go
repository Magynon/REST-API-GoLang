package gateways

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"STORE/domain"

	"github.com/emicklei/go-restful/v3"
)

const elementsPatch = "/elements"

type API struct {
	db map[int]domain.Element
}

func NewAPI() *API {
	return &API{
		db: make(map[int]domain.Element),
	}
}

func (api *API) RegisterRoutes(ws *restful.WebService) {
	ws.Path("/STORE")
	// ws.Route(ws.POST(echoPath).To(api.echoPOSTHandler).Reads(restful.MIME_JSON).Writes(restful.MIME_JSON).Doc("Writes back a json with what you gave it"))
	// ws.Route(ws.GET(echoPath).To(api.echoGETHandler).Writes(restful.MIME_JSON).Doc("Writes back a json with what you gave it"))

	ws.Route(ws.GET(elementsPatch).To(api.GETHandler).Writes(restful.MIME_JSON).Doc("Writes back a json with what you gave it"))
	ws.Route(ws.POST(elementsPatch).To(api.POSTHandler).Writes(restful.MIME_JSON).Doc("Writes back a json with what you gave it"))
	// ws.Route(ws.PATCH(usersPatch).To(api.updateElement).Writes(restful.MIME_JSON).Doc("Writes back a json with what you gave it"))
	// ws.Route(ws.DELETE(usersPatch).To(api.deleteElement).Writes(restful.MIME_JSON).Doc("Writes back a json with what you gave it"))

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
	data, err := ioutil.ReadAll(body)
	if err != nil {
		log.Printf("[ERROR] Couldn't read request body")
		resp.WriteServiceError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()))
		return
	}

	// declare array of products
	var elem = &domain.Element{}
	ok := json.Unmarshal(data, &elem)

	if ok != nil {
		log.Printf("[ERROR] Couldn't unmarshal body data")
	}

	hash := elem.GetHash()

	if _, ok := api.db[hash]; ok  {
		log.Printf("Product already exists!")
		resp.WriteError(http.StatusConflict, fmt.Errorf("Product already exists!"))
		return
	}

	fmt.Println(string(data))
	api.db[hash] = *elem

	resp.Write([]byte("Product successfully added! " + strconv.Itoa(hash)))
}

func (api *API) GETHandler(req *restful.Request, resp *restful.Response) {
	// id_data := req.QueryParameter("id")
	// var id, _ = strconv.Atoi(id_data)

	body := req.Request.Body
	if body == nil {
		log.Printf("[ERROR] Couldn't read request body")
		resp.WriteServiceError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, "nil body"))
		return

	}
	defer body.Close()
	var err error
	data, err := ioutil.ReadAll(body)
	if err != nil {
		log.Printf("[ERROR] Couldn't read request body")
		resp.WriteServiceError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()))
		return
	}

	var id, _ = strconv.Atoi(string(data))
	
	if val, ok := api.db[id]; ok {
		resp.WriteAsJson(val)
	} else {
		log.Printf("[ERROR] Element not found %d", id)
		resp.WriteError(http.StatusNotFound, fmt.Errorf("Element not found"))
		return
	}
}
