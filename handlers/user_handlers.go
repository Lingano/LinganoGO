package handlers

import (
	"LinganoGO/config"
	"LinganoGO/models"
	"LinganoGO/utils"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserProfileResponse defines the structure for the user profile response.
type UserProfileResponse struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Email   string             `json:"email"`
	Profile models.UserProfile `json:"profile"`
}

// UpdateUserProfileRequest defines the structure for updating a user profile.
type UpdateUserProfileRequest struct {
	Bio            *string `json:"bio,omitempty"`
	ProfilePicture *string `json:"profile_picture,omitempty"`
	// Add other updatable profile fields here if needed
}

// AddSavedWordRequest defines the structure for adding a saved word.
type AddSavedWordRequest struct {
	Word       string `json:"word"`
	Definition string `json:"definition,omitempty"`
	Context    string `json:"context,omitempty"`
	ReadingID  string `json:"reading_id,omitempty"`
}

// GetUserProfile godoc
// @Summary      Get user profile
// @Description  Fetches and returns the authenticated user's profile information
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  handlers.UserProfileResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /api/user/profile [get]
// GetUserProfile handles fetching and returning a user's profile.
func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userIDString, ok := r.Context().Value(utils.UserIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Could not retrieve user ID from context"})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(userIDString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user ID format"})
		return
	}

	usersCollection := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err = usersCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch user profile"})
		return
	}

	// Prepare response: do not send sensitive data like password hash
	profileResponse := UserProfileResponse{
		ID:      user.ID.Hex(),
		Name:    user.Name,
		Email:   user.Email,
		Profile: user.Profile,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(profileResponse)
	if err != nil {
		return
	}
}

// UpdateUserProfile godoc
// @Summary      Update user profile
// @Description  Updates the authenticated user's profile information
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        profile  body  handlers.UpdateUserProfileRequest  true  "Profile update info"
// @Success      200  {object}  handlers.UserProfileResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /api/user/profile [put]
// UpdateUserProfile handles updating a user's profile.
func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userIDString, ok := r.Context().Value(utils.UserIDKey).(string) // Changed to utils.UserIDKey
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Could not retrieve user ID from context"})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(userIDString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user ID format"})
		return
	}

	var req UpdateUserProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	usersCollection := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err = usersCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch user for update"})
		return
	}

	// Update fields if provided in the request
	updateFields := bson.M{}
	if req.Bio != nil {
		user.Profile.Bio = *req.Bio
		updateFields["profile.bio"] = *req.Bio
	}
	if req.ProfilePicture != nil {
		user.Profile.ProfilePicture = *req.ProfilePicture
		updateFields["profile.profile_picture"] = *req.ProfilePicture
	}

	if len(updateFields) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "No update fields provided"})
		return
	}

	user.UpdatedAt = time.Now()
	updateFields["updated_at"] = user.UpdatedAt
	// If you are updating nested structs like UserProfile, ensure it's initialized.
	// If user.Profile was nil, this would panic. It's initialized in RegisterUser.
	// user.Profile.LastLogin = time.Now() // Example: update last login, if needed for profile updates
	// updateFields["profile.last_login"] = user.Profile.LastLogin

	_, err = usersCollection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updateFields},
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update user profile"})
		return
	}

	// Prepare response
	updatedProfileResponse := UserProfileResponse{
		ID:      user.ID.Hex(),
		Name:    user.Name,
		Email:   user.Email,
		Profile: user.Profile, // user.Profile now contains the updated fields
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedProfileResponse)
}

// AddSavedWord godoc
// @Summary      Add saved word
// @Description  Adds a new word to the authenticated user's saved words list
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        word  body  handlers.AddSavedWordRequest  true  "Word to save"
// @Success      201  {object}  models.SavedWord
// @Failure      400  {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /api/user/saved-words [post]
// AddSavedWord handles adding a new word to the user's saved words list.
func AddSavedWord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userIDString, ok := r.Context().Value(utils.UserIDKey).(string) // Changed to utils.UserIDKey
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Could not retrieve user ID from context"})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(userIDString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user ID format"})
		return
	}

	var req AddSavedWordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	if req.Word == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Word cannot be empty"})
		return
	}

	newSavedWord := models.SavedWord{
		ID:         primitive.NewObjectID(),
		Word:       req.Word,
		Definition: req.Definition,
		Context:    req.Context,
		DateSaved:  time.Now(),
		ReadingID:  req.ReadingID,
	}

	usersCollection := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$push": bson.M{"saved_words": newSavedWord},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	_, err = usersCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to add saved word"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newSavedWord)
}

// GetSavedWords godoc
// @Summary      Get saved words
// @Description  Fetches and returns the authenticated user's saved words list
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.SavedWord
// @Failure      400  {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /api/user/saved-words [get]
// GetSavedWords handles fetching and returning the user's saved words list.
// SADASD
func GetSavedWords(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userIDString, ok := r.Context().Value(utils.UserIDKey).(string) // Changed to utils.UserIDKey
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Could not retrieve user ID from context"})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(userIDString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user ID format"})
		return
	}

	usersCollection := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err = usersCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch user data"})
		return
	}

	// If SavedWords is nil, return an empty slice to avoid null in JSON response
	savedWords := user.SavedWords
	if savedWords == nil {
		savedWords = []models.SavedWord{}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(savedWords)
}

// DeleteSavedWord godoc
// @Summary      Delete saved word
// @Description  Removes a word from the authenticated user's saved words list
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        savedWordID  path  string  true  "Saved Word ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /api/user/saved-words/{savedWordID} [delete]
// DeleteSavedWord removes a word from the user's saved words list
func DeleteSavedWord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userIDString, ok := r.Context().Value(utils.UserIDKey).(string) // Changed to utils.UserIDKey
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Could not retrieve user ID from context"})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(userIDString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user ID format"})
		return
	}

	vars := mux.Vars(r)
	savedWordIDHex, ok := vars["savedWordID"]
	if !ok {
		http.Error(w, "Saved word ID not provided in path", http.StatusBadRequest)
		return
	}

	savedWordID, err := primitive.ObjectIDFromHex(savedWordIDHex)
	if err != nil {
		http.Error(w, "Invalid saved word ID format", http.StatusBadRequest)
		return
	}

	collection := config.GetCollection("users")
	filter := bson.M{"_id": objectID}
	update := bson.M{"$pull": bson.M{"saved_words": bson.M{"_id": savedWordID}}} // Assuming SavedWord struct has an _id field

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(w, "Failed to delete saved word: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if result.ModifiedCount == 0 {
		// This could mean the user doesn't exist, or the wordID wasn't found.
		// For security, or simplicity, we can return a generic success or a more specific "not found"
		// Checking if the user exists first might be better in a real app.
		http.Error(w, "Word not found or already deleted", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Word deleted successfully"})
}
