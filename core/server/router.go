package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) NewRouter() *mux.Router {
	type Route struct {
		Name        string
		Pattern     string
		Method      string
		HandlerFunc http.HandlerFunc
	}

	routes := []Route{
		Route{
			"UpdateSocket",
			"/updates",
			"GET",
			s.UpdateSocketHandler,
		},
		Route{
			"GroupIndex",
			"/groups",
			"GET",
			s.GroupIndexHandler,
		},
		Route{
			"Group",
			"/groups/{id}",
			"GET",
			s.GroupHandler,
		},
		Route{
			"GroupCreate",
			"/groups",
			"POST",
			s.GroupCreateHandler,
		},
		Route{
			"GroupExport",
			"/groups/{id}/export",
			"GET",
			s.GroupExportHandler,
		},
		Route{
			"GroupImport",
			"/groups/{id}/import",
			"POST",
			s.GroupImportHandler,
		},
		Route{
			"GroupModifyLabel",
			"/groups/{id}/label",
			"PUT",
			s.GroupModifyLabelHandler,
		},
		Route{
			"GroupModifyAllChildren",
			"/groups/{id}/children",
			"PUT",
			s.GroupModifyAllChildrenHandler,
		},
		Route{
			"GroupModifyChild",
			"/groups/{id}/children/{node_id}",
			"PUT",
			s.GroupModifyChildHandler,
		},
		Route{
			"GroupPosition",
			"/groups/{id}/position",
			"PUT",
			s.GroupPositionHandler,
		},
		Route{
			"GroupDelete",
			"/groups/{id}",
			"DELETE",
			s.GroupDeleteHandler,
		},
		Route{
			"BlockIndex",
			"/blocks",
			"GET",
			s.BlockIndexHandler,
		},
		Route{
			"Block",
			"/blocks/{id}",
			"GET",
			s.BlockHandler,
		},
		Route{
			"BlockCreate",
			"/blocks",
			"POST",
			s.BlockCreateHandler,
		},
		Route{
			"BlockDelete",
			"/blocks/{id}",
			"DELETE",
			s.BlockDeleteHandler,
		},
		Route{
			"BlockModifyName",
			"/blocks/{id}/label",
			"PUT",
			s.BlockModifyNameHandler,
		},
		Route{
			"BlockModifyRoute",
			"/blocks/{id}/routes/{index}",
			"PUT",
			s.BlockModifyRouteHandler,
		},
		Route{
			"BlockModifyPosition",
			"/blocks/{id}/position",
			"PUT",
			s.BlockModifyPositionHandler,
		},
		Route{
			"ConnectionIndex",
			"/connections",
			"GET",
			s.ConnectionIndexHandler,
		},
		Route{
			"ConnectionCreate",
			"/connections",
			"POST",
			s.ConnectionCreateHandler,
		},
		Route{
			"ConnectionModifyCoordinates",
			"/connections/{id}/coordinates",
			"PUT",
			s.ConnectionModifyCoordinates,
		},
		Route{
			"ConnectionDelete",
			"/connections/{id}",
			"DELETE",
			s.ConnectionDeleteHandler,
		},
		Route{
			"SourceCreate",
			"/sources",
			"POST",
			s.SourceCreateHandler,
		},
		Route{
			"SourceIndex",
			"/sources",
			"GET",
			s.SourceIndexHandler,
		},
		Route{
			"SourceModify",
			"/sources/{id}",
			"PUT",
			s.SourceModifyHandler,
		},
		Route{
			"Source",
			"/sources/{id}",
			"GET",
			s.SourceHandler,
		},
		Route{
			"Source",
			"/sources/{id}",
			"DELETE",
			s.SourceDeleteHandler,
		},
		Route{
			"Library",
			"/library",
			"GET",
			s.LibraryHandler,
		},
		Route{
			"LinkIndex",
			"/links",
			"GET",
			s.LinkIndexHandler,
		},
		Route{
			"LinkCreate",
			"/links",
			"POST",
			s.LinkCreateHandler,
		},
		Route{
			"LinkDelete",
			"/connections/{id}",
			"DELETE",
			s.ConnectionDeleteHandler,
		},
	}
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	return router

}