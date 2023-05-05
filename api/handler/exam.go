package handler

import (
	"app/api/models"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) EachStaff(c *gin.Context) {
	year := c.Param("year")

	fmt.Println(year)

	resp, err := h.storages.Exam().EachStaff(context.Background(), &models.Date{Day: year})
	if err != nil {
		h.handlerResponse(c, "storage.Exam.StaffDate", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "all staff date", http.StatusCreated, resp)

}

func (h *Handler) Total(c *gin.Context) {
	var data models.Id

	err := c.ShouldBindJSON(&data)
	if err != nil {
		h.handlerResponse(c, "error", http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.storages.Exam().Total(context.Background(), &data)
	if err != nil {
		h.handlerResponse(c, "server error", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "res", http.StatusCreated, resp)
}

func (h *Handler) Create(c *gin.Context) {

	var createPromoCode models.CreatePromoCode

	err := c.ShouldBindJSON(&createPromoCode) // parse req body to given type struct
	if err != nil {
		h.handlerResponse(c, "create promocode", http.StatusBadRequest, err.Error())
		return
	}
	// fmt.Println(createPromoCode)
	id, err := h.storages.Exam().Create(context.Background(), &createPromoCode)
	if err != nil {
		h.handlerResponse(c, "storage.order.create", http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := h.storages.Exam().GetByID(context.Background(), &models.PromocodePrimaryKey{PromocodeId: id})
	if err != nil {
		h.handlerResponse(c, "storage.promocode.getByID", http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println(resp)
	h.handlerResponse(c, "create order", http.StatusCreated, id)
}

func (h *Handler) GetByIdPromocode(c *gin.Context) {

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		h.handlerResponse(c, "storage.order.getByID", http.StatusBadRequest, "id incorrect")
		return
	}

	resp, err := h.storages.Exam().GetByID(context.Background(), &models.PromocodePrimaryKey{PromocodeId: idInt})
	if err != nil {
		h.handlerResponse(c, "storage.order.getByID", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get order by id", http.StatusCreated, resp)
}

func (h *Handler) GetListPromocode(c *gin.Context) {

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.handlerResponse(c, "get list order", http.StatusBadRequest, "invalid offset")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.handlerResponse(c, "get list order", http.StatusBadRequest, "invalid limit")
		return
	}

	resp, err := h.storages.Exam().GetList(context.Background(), &models.GetListBrandRequest{
		Offset: offset,
		Limit:  limit,
		Search: c.Query("search"),
	})
	if err != nil {
		h.handlerResponse(c, "storage.order.getlist", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get list order response", http.StatusOK, resp)
}

func (h *Handler) DeletePromocode(c *gin.Context) {

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		h.handlerResponse(c, "storage.promocode.getByID", http.StatusBadRequest, "id incorrect")
		return
	}

	rowsAffected, err := h.storages.Exam().Delete(context.Background(), &models.PromocodePrimaryKey{PromocodeId: idInt})
	if err != nil {
		h.handlerResponse(c, "storage.promocode.delete", http.StatusInternalServerError, err.Error())
		return
	}
	if rowsAffected <= 0 {
		h.handlerResponse(c, "storage.promocode.delete", http.StatusBadRequest, "now rows affected")
		return
	}

	h.handlerResponse(c, "delete order", http.StatusNoContent, nil)
}
