package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Zorynix/song-library/internal/entity"
	errs "github.com/Zorynix/song-library/internal/errors"
	logger "github.com/Zorynix/song-library/internal/logger"
	"github.com/Zorynix/song-library/internal/services"
	"github.com/go-chi/chi/v5"
)

// @title Song Library API
// @version 1.0
// @description API для управления библиотекой песен
// @host localhost:8080
// @BasePath /api/v1

type Handler struct {
	services *services.Services
}

func NewHandler(services *services.Services) *Handler {
	return &Handler{services: services}
}

func (h *Handler) Register(r chi.Router) {
	r.Get("/songs", h.GetSongs)
	r.Get("/songs/{id}/verses", h.GetSongVerses)
	r.Delete("/songs/{id}", h.DeleteSong)
	r.Put("/songs/{id}", h.UpdateSong)
	r.Post("/songs", h.AddSong)
}

// GetSongs возвращает список песен с фильтрацией
// @Summary Получить список песен
// @Description Возвращает список песен с возможностью фильтрации по группе, названию и тексту
// @Tags Songs
// @Accept json
// @Produce json
// @Param group query string false "Название группы"
// @Param song query string false "Название песни"
// @Param text query string false "Текст песни"
// @Param limit query int false "Лимит записей"
// @Param offset query int false "Смещение"
// @Success 200 {array} entity.Song "Список песен"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /songs [get]
func (h *Handler) GetSongs(w http.ResponseWriter, r *http.Request) {
	var filter entity.SongFilter
	filter.Group = r.URL.Query().Get("group")
	filter.Title = r.URL.Query().Get("song")
	filter.Text = r.URL.Query().Get("text")
	if limit := r.URL.Query().Get("limit"); limit != "" {
		filter.Limit, _ = strconv.Atoi(limit)
	}
	if offset := r.URL.Query().Get("offset"); offset != "" {
		filter.Offset, _ = strconv.Atoi(offset)
	}

	logger.Logger.Debug().
		Str("group", filter.Group).
		Str("title", filter.Title).
		Str("text", filter.Text).
		Int("limit", filter.Limit).
		Int("offset", filter.Offset).
		Msg("Handling GetSongs request")

	songs, err := h.services.Song.GetSongs(r.Context(), filter)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("Failed to handle GetSongs request")
		http.Error(w, errs.ErrInternal.Error(), http.StatusInternalServerError)
		return
	}

	logger.Logger.Info().Int("count", len(songs)).Msg("GetSongs request handled successfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

// GetSongVerses возвращает куплеты песни по ID
// @Summary Получить куплеты песни
// @Description Возвращает куплеты песни по её ID с пагинацией
// @Tags Songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param limit query int false "Лимит куплетов"
// @Param offset query int false "Смещение"
// @Success 200 {array} string "Список куплетов"
// @Failure 400 {string} string "Неверный ID"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /songs/{id}/verses [get]
func (h *Handler) GetSongVerses(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	var pagination entity.VersePagination
	pagination.SongID = id
	if limit := r.URL.Query().Get("limit"); limit != "" {
		pagination.Limit, _ = strconv.Atoi(limit)
	}
	if offset := r.URL.Query().Get("offset"); offset != "" {
		pagination.Offset, _ = strconv.Atoi(offset)
	}

	logger.Logger.Debug().
		Int64("song_id", pagination.SongID).
		Int("limit", pagination.Limit).
		Int("offset", pagination.Offset).
		Msg("Handling GetSongVerses request")

	verses, err := h.services.Song.GetSongVerses(r.Context(), pagination)
	if err != nil {
		logger.Logger.Error().Err(err).Int64("song_id", id).Msg("Failed to handle GetSongVerses request")
		if errors.Is(err, errs.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, errs.ErrInternal.Error(), http.StatusInternalServerError)
		return
	}

	logger.Logger.Info().
		Int64("song_id", id).
		Int("verse_count", len(verses)).
		Msg("GetSongVerses request handled successfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(verses)
}

// DeleteSong удаляет песню по ID
// @Summary Удалить песню
// @Description Удаляет песню по её ID
// @Tags Songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Success 204 {string} string "Песня успешно удалена"
// @Failure 400 {string} string "Неверный ID"
// @Failure 404 {string} string "Песня не найдена"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /songs/{id} [delete]
func (h *Handler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	logger.Logger.Debug().Int64("id", id).Msg("Handling DeleteSong request")

	err := h.services.Song.DeleteSong(r.Context(), id)
	if err != nil {
		logger.Logger.Error().Err(err).Int64("id", id).Msg("Failed to handle DeleteSong request")
		if errors.Is(err, errs.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, errs.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, errs.ErrInternal.Error(), http.StatusInternalServerError)
		return
	}

	logger.Logger.Info().Int64("id", id).Msg("DeleteSong request handled successfully")
	w.WriteHeader(http.StatusNoContent)
}

// UpdateSong обновляет песню по ID
// @Summary Обновить песню
// @Description Обновляет данные песни по её ID
// @Tags Songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param song body entity.Song true "Данные песни"
// @Success 200 {object} entity.Song "Обновленная песня"
// @Failure 400 {string} string "Неверный запрос или ID"
// @Failure 404 {string} string "Песня не найдена"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /songs/{id} [put]
func (h *Handler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	var song entity.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		logger.Logger.Error().Err(err).Msg("Failed to decode song data")
		http.Error(w, errs.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}
	song.ID = id

	logger.Logger.Debug().
		Int64("id", song.ID).
		Str("group", song.Group).
		Str("title", song.Title).
		Msg("Handling UpdateSong request")

	err := h.services.Song.UpdateSong(r.Context(), song)
	if err != nil {
		logger.Logger.Error().Err(err).Int64("id", id).Msg("Failed to handle UpdateSong request")
		if errors.Is(err, errs.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, errs.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, errs.ErrInternal.Error(), http.StatusInternalServerError)
		return
	}

	logger.Logger.Info().Int64("id", id).Msg("UpdateSong request handled successfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(song)
}

// AddSong добавляет новую песню
// @Summary Добавить песню
// @Description Добавляет новую песню, обогащая её данными из внешнего API
// @Tags Songs
// @Accept json
// @Produce json
// @Param song body entity.Song true "Данные песни (group и title обязательны)"
// @Success 201 {object} entity.Song "Созданная песня"
// @Failure 400 {string} string "Неверный запрос"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /songs [post]
func (h *Handler) AddSong(w http.ResponseWriter, r *http.Request) {
	var song entity.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		logger.Logger.Error().Err(err).Msg("Failed to decode song data")
		http.Error(w, errs.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}

	logger.Logger.Debug().
		Str("group", song.Group).
		Str("title", song.Title).
		Msg("Handling AddSong request")

	createdSong, err := h.services.Song.AddSong(r.Context(), song)
	if err != nil {
		logger.Logger.Error().Err(err).Str("group", song.Group).Str("title", song.Title).Msg("Failed to handle AddSong request")
		http.Error(w, errs.ErrInternal.Error(), http.StatusInternalServerError)
		return
	}

	logger.Logger.Info().Int64("id", createdSong.ID).Msg("AddSong request handled successfully")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdSong)
}
