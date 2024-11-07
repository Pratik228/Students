package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func isValidEmail (email string) bool {
	email = strings.TrimSpace(strings.ToLower(email))
	return emailRegex.MatchString(email)
}

func validateStudent(student Student) error {
    if student.StudentID == 0 {
        return fmt.Errorf("Student ID is required")
    }
    if student.Name == "" {
        return fmt.Errorf("Student Name is required")
    }
    if student.Age == 0 {
        return fmt.Errorf("Student Age is required")
    }
    if student.Email == "" {
        return fmt.Errorf("Student Email is required")
    }
    if !isValidEmail(student.Email) {
        return fmt.Errorf("Student Email is invalid")
    }
    return nil
}

type Student struct {
	StudentID int `json:"id"`
	Name string `json:"name"`
	Age int `json:"age"`
	Email string `json:"email"`
}

var  (
	students = make(map[int]Student)
	mutex sync.RWMutex
)

func GetAllStudents(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	mutex.RLock();
	defer mutex.RUnlock();

	studentList := make([]Student, 0, len(students))

	for _, student := range students {
		studentList = append(studentList, student)
	}

	json.NewEncoder(w).Encode(studentList);

}

func CreateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json");

	var newStudent Student;

	if err:= json.NewDecoder(r.Body).Decode(&newStudent); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	if err := validateStudent(newStudent); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	mutex.Lock()
	defer mutex.Unlock()

	if _, exists := students[int(newStudent.StudentID)]; exists {
		http.Error(w, "Student already exists", http.StatusConflict)
		return
	}
	students[int(newStudent.StudentID)] = newStudent
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newStudent)


}

func GetStudentById(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err!=nil {
		http.Error(w, "Invalid student Id format", http.StatusBadRequest)
		return
	}
	
	mutex.RLock();
	defer mutex.RUnlock();

	student, exists := students[id]

	if !exists {
		http.Error(w, "Student does not exist", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(student)

}

func UpdateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id");

	id, err := strconv.Atoi(idStr)

	if err !=nil {
		http.Error(w, "Invalid student Id format", http.StatusBadRequest)
		return
	}

	var updatedStudent Student
	if err:= json.NewDecoder(r.Body).Decode(&updatedStudent); err !=nil {
		http.Error(w, "Invalid request Body", http.StatusBadRequest)
		return
	}

	if err := validateStudent(updatedStudent); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	mutex.Lock();
	defer mutex.Unlock();

	_, exists := students[id]

	if !exists {
		http.Error(w, "Student does not exist", http.StatusBadRequest);
		return
	}

	updatedStudent.StudentID = id
	students[id] = updatedStudent

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedStudent)
}

func DeleteStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id");
	id, err := strconv.Atoi(idStr)

	if err!=nil {
		http.Error(w, "Invalid Student Id ", http.StatusBadRequest)
	}

	mutex.Lock();
	defer mutex.Unlock();

	_, exists := students[id];

	if !exists {
		http.Error(w, "Student does not exist", http.StatusBadRequest)
	}

	delete(students, id);

	w.WriteHeader(http.StatusNoContent)
}

func GetStudentSummary(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid student ID format", http.StatusBadRequest)
        return
    }

    mutex.RLock()
    student, exists := students[id]
    mutex.RUnlock()

    if !exists {
        http.Error(w, "Student not found", http.StatusNotFound)
        return
    }

    prompt := fmt.Sprintf(`Create a professional student profile summary for:
        Student ID: %d
        Name: %s
        Age: %d
        Email: %s

        Write 2-3 sentences describing the student, including their ID, age, and potential academic interests.
        Keep it professional and positive.`,
        student.StudentID,
        student.Name,
        student.Age,
        student.Email,
    )

    requestBody := map[string]interface{}{
        "model":       "mistral", 
        "prompt":      prompt,
        "temperature": 0.7,
        "max_tokens":  300,
        "stream":      false, 
    }

    requestJSON, err := json.Marshal(requestBody)
    if err != nil {
        http.Error(w, "Failed to create summary request", http.StatusInternalServerError)
        return
    }

    fmt.Printf("Sending request to Ollama: %s\n", string(requestJSON))

    resp, err := http.Post(
        "http://localhost:11434/api/generate",
        "application/json",
        bytes.NewBuffer(requestJSON),
    )
    if err != nil {
        http.Error(w, "Failed to generate summary", http.StatusInternalServerError)
        fmt.Printf("Error making request: %v\n", err)
        return
    }
    defer resp.Body.Close()
    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        http.Error(w, "Failed to read response", http.StatusInternalServerError)
        return
    }
    fmt.Printf("Raw Ollama response: %s\n", string(respBody))

    var ollamaResponse struct {
        Response string `json:"response"`
    }
    if err := json.Unmarshal(respBody, &ollamaResponse); err != nil {
        http.Error(w, "Failed to parse summary", http.StatusInternalServerError)
        fmt.Printf("Error parsing response: %v\n", err)
        return
    }

    if len(ollamaResponse.Response) < 10 {
        http.Error(w, "Generated summary too short", http.StatusInternalServerError)
        return
    }

    summary := struct {
        StudentID int    `json:"student_id"`
        Name      string `json:"name"`
        Summary   string `json:"summary"`
    }{
        StudentID: student.StudentID,
        Name:      student.Name,
        Summary:   strings.TrimSpace(ollamaResponse.Response),
    }

    json.NewEncoder(w).Encode(summary)
}