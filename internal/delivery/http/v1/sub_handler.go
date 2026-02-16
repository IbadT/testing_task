package v1

import (
	"context"
	myerrors "testingtask/internal/errors"
	"testingtask/internal/service"
	"testingtask/internal/web/subscriptions"
	logger "testingtask/pkg"

	"github.com/google/uuid"
)

type SubHandler struct {
	serv service.SubService
}

func NewSubHandler(s service.SubService) *SubHandler {
	return &SubHandler{serv: s}
}

// Create Создать подписку
// @Summary Создать подписку
// @Description Создаёт новую подписку и возвращает её идентификатор
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body SubscriptionDTO true "Данные для создания подписки"
// @Success 201 {object} SubscriptionID "Подписка успешно создана"
// @Failure 400 {object} myerrors.ErrorResponse "Некорректные данные"
// @Failure 500 {object} myerrors.ErrorInternalServerError "Внутренняя ошибка сервера"
// @Router /subscriptions [post]
func (h *SubHandler) Create(ctx context.Context, request subscriptions.CreateRequestObject) (subscriptions.CreateResponseObject, error) {
	logger.Info(ctx, "create subscription called", map[string]interface{}{
		"body": request.Body,
	})
	dto := CreateRequestToDTO(*request.Body)

	domainObj, err := DTOToDomain(nil, *dto)
	if err != nil {
		logger.Error(ctx, "invalid data", err, nil)
		resp, _ := myerrors.MapError(myerrors.ErrInvalidData)
		return subscriptions.Create400JSONResponse(resp), nil
	}

	id, err := h.serv.Create(ctx, domainObj)
	if err != nil {
		logger.Error(ctx, "error create subscripton", err, nil)
		resp, _ := myerrors.MapError(err)
		return subscriptions.Create500JSONResponse(resp), nil
	}

	return CreateToResponse(id), nil
}

// Get Получить подписку по ID
// @Summary Получить подписку по ID
// @Description Возвращает подписку по её идентификатору
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Success 200 {object} SubscriptionResponseDTO "Подписка найдена"
// @Failure 400 {object} myerrors.ErrorResponse "Некорректный ID"
// @Failure 404 {object} myerrors.ErrorNotFound "Подписка не найдена"
// @Failure 500 {object} myerrors.ErrorInternalServerError "Внутренняя ошибка сервера"
// @Router /subscriptions/{id} [get]
func (h *SubHandler) Get(ctx context.Context, request subscriptions.GetRequestObject) (subscriptions.GetResponseObject, error) {
	logger.Info(ctx, "update subscription called", map[string]interface{}{
		"id": request.Id,
	})
	uid, err := uuid.Parse(request.Id)
	if err != nil {
		logger.Error(ctx, "invalid id format", err, nil)
		resp, _ := myerrors.MapError(myerrors.ErrInvalidID)
		return subscriptions.Get400JSONResponse(resp), nil
	}

	subscription, err := h.serv.Get(ctx, uid)
	if err != nil {
		logger.Error(ctx, "subscripton not found", err, nil)
		resp, code := myerrors.MapError(err)
		switch code {
		case 404:
			return subscriptions.Get404JSONResponse(resp), nil
		default:
			return subscriptions.Get500JSONResponse(resp), nil
		}
	}

	resDto := DomainToDTO(subscription)

	return GetDTOToResponse(resDto), nil
}

// List Получить список подписок
// @Summary Получить список подписок
// @Description Возвращает список подписок с пагинацией
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param limit query int true "Количество элементов"
// @Param offset query int true "Смещение"
// @Success 200 {array} SubscriptionResponseDTO "Список подписок"
// @Failure 400 {object} myerrors.ErrorResponse "Некорректный ID"
// @Failure 404 {object} myerrors.ErrorNotFound "Подписка не найдена"
// @Failure 500 {object} myerrors.ErrorInternalServerError "Внутренняя ошибка сервера"
// @Router /subscriptions [get]
func (h *SubHandler) List(ctx context.Context, request subscriptions.ListRequestObject) (subscriptions.ListResponseObject, error) {
	logger.Info(ctx, "list subscriptions called", map[string]interface{}{
		"params": request.Params,
	})
	dto := ListRequestToDTO(request)

	paging := NewPagingBase(dto)

	subs, totalCount, err := h.serv.List(ctx, paging)
	if err != nil {
		logger.Error(ctx, "error list", err, nil)
		resp, code := myerrors.MapError(err)

		switch code {
		case 400:
			return subscriptions.List400JSONResponse(resp), nil
		case 404:
			return subscriptions.List404JSONResponse(resp), nil
		default:
			return subscriptions.List500JSONResponse(resp), nil
		}
	}

	resDto := DomainListToDTO(subs)

	return ListDTOToResponse(resDto, dto, totalCount), nil
}

// Sum Получить сумму стоимости подписок
// @Summary Получить сумму стоимости подписок
// @Description Возвращает суммарную стоимость всех подписок пользователя
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "User ID"
// @Param service_name query string false "Service name"
// @Param start query string false "Start date (MM-YYYY)"
// @Param end query string false "End date (MM-YYYY)"
// @Param limit query int false "Limit subscriptions for count price" default(10)
// @Param offset query int false "Offset subscriptions" default(0)
// @Success 200 {object} ListSubscriptionsResponseDto
// @Failure 400 {object} myerrors.ErrorResponse "Некорректный ID"
// @Failure 404 {object} myerrors.ErrorNotFound "Подписка не найдена"
// @Failure 500 {object} myerrors.ErrorInternalServerError "Внутренняя ошибка сервера"
// @Router /subscriptions/sum [get]
func (h *SubHandler) Sum(ctx context.Context, request subscriptions.SumRequestObject) (subscriptions.SumResponseObject, error) {
	logger.Info(ctx, "sum subscriptions called", map[string]interface{}{
		"params": request.Params,
	})
	dto := SumRequestToDTO(request)

	filter, err := SumDTOToDomain(dto)
	if err != nil {
		logger.Error(ctx, "invalid data", err, nil)
		resp, code := myerrors.MapError(err)
		switch code {
		case 400:
			return subscriptions.Sum400JSONResponse(resp), nil
		default:
			return subscriptions.Sum500JSONResponse(resp), nil
		}
	}

	result, err := h.serv.Sum(ctx, filter)
	if err != nil {
		logger.Error(ctx, "error sum subscriptions", err, nil)
		resp, code := myerrors.MapError(err)
		switch code {
		case 400:
			return subscriptions.Sum400JSONResponse(resp), nil
		case 404:
			return subscriptions.Sum404JSONResponse(resp), nil
		default:
			return subscriptions.Sum500JSONResponse(resp), nil
		}
	}

	paging := NewPagingDTO(filter.Limit, filter.Offset, result.TotalCount)
	responseDTO := DomainToSumDTO(result)

	return SumDTOToResponse(responseDTO, paging), nil
}

// Update Обновить подписку
// @Summary Обновить подписку
// @Description Обновляет данные подписки по ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Param request body SubscriptionDTO true "Данные для обновления подписки"
// @Success 200 {object} SubscriptionID "Подписка успешно обновлена"
// @Failure 400 {object} myerrors.ErrorResponse "Некорректные данные"
// @Failure 404 {object} myerrors.ErrorNotFound "Подписка не найдена"
// @Failure 500 {object} myerrors.ErrorInternalServerError "Внутренняя ошибка сервера"
// @Router /subscriptions/{id} [put]
func (h *SubHandler) Update(ctx context.Context, request subscriptions.UpdateRequestObject) (subscriptions.UpdateResponseObject, error) {
	logger.Info(ctx, "udpate subscriptions called", map[string]interface{}{
		"id":   request.Id,
		"body": request.Body,
	})
	uid, err := uuid.Parse(request.Id.String())
	if err != nil {
		logger.Error(ctx, "invalid id format", err, nil)
		resp, code := myerrors.MapError(myerrors.ErrInvalidID)
		switch code {
		case 400:
			return subscriptions.Update400JSONResponse(resp), nil
		default:
			return subscriptions.Update500JSONResponse(resp), nil
		}
	}

	dto := UpdateRequestToDTO(*request.Body)

	subDomain, err := DTOToDomain(&uid, *dto)
	if err != nil {
		logger.Error(ctx, "invalid data", err, nil)
		resp, code := myerrors.MapError(err)
		switch code {
		case 400:
			return subscriptions.Update400JSONResponse(resp), nil
		default:
			return subscriptions.Update500JSONResponse(resp), nil
		}
	}

	id, err := h.serv.Update(ctx, uid, subDomain)
	if err != nil {
		logger.Error(ctx, "error update subscription", err, nil)
		resp, code := myerrors.MapError(err)
		switch code {
		case 400:
			return subscriptions.Update400JSONResponse(resp), nil
		case 404:
			return subscriptions.Update404JSONResponse(resp), nil
		default:
			return subscriptions.Update500JSONResponse(resp), nil
		}
	}

	return UpdateToResponse(id), nil
}

// Delete Удалить подписку
// @Summary Удалить подписку
// @Description Удаляет подписку по её идентификатору
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Success 204 "Подписка успешно удалена"
// @Failure 400 {object} myerrors.ErrorResponse "Некорректный ID"
// @Failure 404 {object} myerrors.ErrorNotFound "Подписка не найдена"
// @Failure 500 {object} myerrors.ErrorInternalServerError "Внутренняя ошибка сервера"
// @Router /subscriptions/{id} [delete]
func (h *SubHandler) Delete(ctx context.Context, request subscriptions.DeleteRequestObject) (subscriptions.DeleteResponseObject, error) {
	logger.Info(ctx, "delete subscription called", map[string]interface{}{
		"id": request.Id,
	})
	uid, err := uuid.Parse(request.Id.String())
	if err != nil {
		logger.Error(ctx, "invalid id format", err, nil)
		resp, _ := myerrors.MapError(myerrors.ErrInvalidID)
		return subscriptions.Delete400JSONResponse(resp), nil
	}

	if err = h.serv.Delete(ctx, uid); err != nil {
		logger.Error(ctx, "error delete subscription", err, nil)
		resp, code := myerrors.MapError(err)
		switch code {
		case 404:
			return subscriptions.Delete404JSONResponse(resp), nil
		default:
			return subscriptions.Delete500JSONResponse(resp), nil
		}
	}

	return subscriptions.Delete204Response{}, nil
}
