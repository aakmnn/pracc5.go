package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "strings"

    "practice5_ready/internal/models"
    "practice5_ready/internal/repository"
)

type Handler struct {
    repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
    return &Handler{repo: repo}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
        return
    }

    writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
        return
    }

    query := r.URL.Query()

    limit, err := parseIntWithDefault(query.Get("limit"), 5)
    if err != nil || limit <= 0 {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "limit must be a positive integer"})
        return
    }

    offset, err := parseIntWithDefault(query.Get("offset"), 0)
    if err != nil || offset < 0 {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "offset must be zero or a positive integer"})
        return
    }

    filters := models.UserFilters{
        Name:      strings.TrimSpace(query.Get("name")),
        Email:     strings.TrimSpace(query.Get("email")),
        Gender:    strings.TrimSpace(strings.ToLower(query.Get("gender"))),
        BirthDate: strings.TrimSpace(query.Get("birth_date")),
        Limit:     limit,
        Offset:    offset,
        OrderBy:   strings.TrimSpace(query.Get("order_by")),
    }

    if idRaw := strings.TrimSpace(query.Get("id")); idRaw != "" {
        id, err := strconv.Atoi(idRaw)
        if err != nil {
            writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id must be an integer"})
            return
        }
        filters.ID = &id
    }

    response, err := h.repo.GetPaginatedUsers(filters)
    if err != nil {
        writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
        return
    }

    writeJSON(w, http.StatusOK, response)
}

func (h *Handler) GetCommonFriends(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
        return
    }

    query := r.URL.Query()
    user1, err1 := strconv.Atoi(strings.TrimSpace(query.Get("user1")))
    user2, err2 := strconv.Atoi(strings.TrimSpace(query.Get("user2")))
    if err1 != nil || err2 != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "user1 and user2 must be integers"})
        return
    }
    if user1 == user2 {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "user1 and user2 must be different"})
        return
    }

    friends, err := h.repo.GetCommonFriends(user1, user2)
    if err != nil {
        writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
        return
    }

    writeJSON(w, http.StatusOK, map[string]any{
        "user1":         user1,
        "user2":         user2,
        "common_friends": friends,
        "count":         len(friends),
    })
}

func parseIntWithDefault(value string, fallback int) (int, error) {
    if strings.TrimSpace(value) == "" {
        return fallback, nil
    }
    return strconv.Atoi(value)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(data)
}
