package handlers

import (
	"../libs"
	"../models"
	"../utils"
	"strconv"
)

type MainHandler struct {
	libs.BaseHandler
}

func (self *MainHandler) Get() {

	inputs := self.Input()
	page, _ := strconv.Atoi(inputs.Get("page"))
	id, _ := strconv.Atoi(self.Ctx.Params[":cid"])
	rcs := len(models.GetAllNodeByCid(id, 0, 0, "hotness"))

	limit := 25
	pages, pageout, beginnum, endnum, offset := utils.Pages(rcs, page, limit)
	self.Data["pagesbar"] = utils.Pagesbar("", rcs, pages, pageout, beginnum, endnum, 1)
	self.Data["nodes_hotness"] = models.GetAllNodeByCid(id, offset, limit, "hotness")
	self.Layout = "layout.html"
	self.TplNames = "index.html"
	self.Render()

}
