package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/mail"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// hold all the data
const dataFile = "hotel_data.json"

// Room represents a hotel room
type Room struct {
	Number     int     `json:"number"`
	Type       string  `json:"type"`
	Price      float64 `json:"price"`
	IsReserved bool    `json:"is_reserved"`
}

// reservation represents a booking
type Reservation struct {
	ID         int       `json:"id"`
	RoomNumber int       `json:"room_number"`
	GuestName  string    `json:"guest_name"`
	GuestEmail string    `json:"guest_email"`
	CheckIn    time.Time `json:"check_in"`
	CheckOut   time.Time `json:"check_out"`
	CreatedAt  time.Time `json:"created_at"`
}

// manages the data to persist
type HotelData struct {
	Rooms        []Room        `json:"rooms"`
	Reservations []Reservation `json:"reservations"`
	NextResID    int           `json:"next_res_id"`
}

// HotelSystem manages the hotel operations
type HotelSystem struct {
	rooms        []Room
	reservations []Reservation
	nextResID    int
}

// hotel system manages rooom, reservation and reservation ID
func NewHotelSystem() *HotelSystem {
	return &HotelSystem{
		rooms: []Room{
			{Number: 101, Type: "Single", Price: 80.00, IsReserved: false},
			{Number: 102, Type: "Double", Price: 100.00, IsReserved: false},
			{Number: 201, Type: "Suite", Price: 280.00, IsReserved: false},
			{Number: 103, Type: "Single", Price: 90.00, IsReserved: false},
			{Number: 202, Type: "Double", Price: 130.00, IsReserved: false},
			{Number: 301, Type: "Suite", Price: 289.00, IsReserved: false},
		},

		reservations: []Reservation{},
		nextResID:    1,
	}
}

// loadData loads hotel datafrom json file
func (h *HotelSystem) loadData() error {
	file, err := os.Open(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // file doesnt exist yet, use the defaults
		}

		return err
	}

	defer file.Close()

	var data HotelData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return err
	}

	h.rooms = data.Rooms
	h.reservations = data.Reservations
	h.nextResID = data.NextResID

	fmt.Println("Data loaded successfully from ", dataFile)
	return nil
}

// saveData saves hotel data to Json
func (h *HotelSystem) saveData() error {
	data := HotelData{
		Rooms:        h.rooms,
		Reservations: h.reservations,
		NextResID:    h.nextResID,
	}

	file, err := os.Create(dataFile)
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}

// the main function, everything starts here
func main() {
	h := NewHotelSystem()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("======================================================")
	fmt.Println("      Welcome to the new Hotel Revervation System")
	fmt.Println("======================================================")

	for {
		displayMenu()
		option := readInput(reader, "Select an option: ")

		switch option {
		case "1":
			h.viewAvailableRooms()
		case "2":
			h.makeReservation(reader)
		case "3":
			h.viewReservations()
		case "4":
			h.cancelReservations(reader)
		case "5":
			h.searchRoom(reader)
		case "6":
			fmt.Println("\nThank you for using our system. Goodbye...")
			return
		default:
			fmt.Println("Invalid option selected. Please select 1-6.")
		}
		fmt.Println()
	}
}

// display menu for all functions
func displayMenu() {
	fmt.Println("--------------------------------------")
	fmt.Println("|      MAIN MENU                     |")
	fmt.Println("--------------------------------------")
	fmt.Println("1. View Available Rooms")
	fmt.Println("2. Make Reservation")
	fmt.Println("3. View Reservations")
	fmt.Println("4. Cancel Reservations")
	fmt.Println("5. Search Room by Type")
	fmt.Println("6. Exit")
	fmt.Println("--------------------------------------- ")

}

// reads inputs and removes spaces or newlines
func readInput(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// fucntion for viewing all rooms
func (h *HotelSystem) viewAvailableRooms() {
	fmt.Println("\n Available Rooms:")
	fmt.Println("----------------------------------------")

	hasAvailable := false
	for _, room := range h.rooms {
		if !room.IsReserved {
			fmt.Printf(" Room %d | Type: %-8s | Price: $%.2f/night\n",
				room.Number, room.Type, room.Price)
			hasAvailable = true

		}
	}

	if !hasAvailable {
		fmt.Println("No Romms currently available")
	}
	fmt.Println("--------------------------------------")
}

// func for making room reservations
func (h *HotelSystem) makeReservation(reader *bufio.Reader) {
	h.viewAvailableRooms()

	// get room number
	roomInput := readInput(reader, "\nEnter Room Number to reserve:")
	roomInput = strings.TrimSpace(roomInput)
	roomNumber, err := strconv.Atoi(roomInput)
	if err != nil {
		fmt.Println("Invalid room number. Please enter a valid number")
		return
	}

	// find and validate room
	roomIndex := h.findRoomIndex(roomNumber)
	if roomIndex == -1 {
		fmt.Println("Room not found")
		return
	}

	if h.rooms[roomIndex].IsReserved {
		fmt.Println("Room is already reserved")
		return
	}

	// get Guest information
	guestName := readInput(reader, "Enter Guest Name: ")
	if guestName == "" {
		fmt.Println("Guest name cannot be empty.")
		return
	}

	// gets and validates guest email
	guestEmail := readInput(reader, "Enter Guest Email: ")
	if !isValidEmail(guestEmail) {
		fmt.Println("Warrning: Email format may be invalid.")
		return
	}

	// get check in date
	checkInStr := readInput(reader, "Enter Check-In date (YYYY-MM-DD) :")
	checkIn, err := time.Parse("2006-01-02", checkInStr)
	if err != nil {
		fmt.Println("Invalid date format. Use YYYY-MM-DD")
		return
	}

	// get check out date
	checkOutStr := readInput(reader, "Enter Check-Out date (YYYY-MM-DD) :")
	checkOut, err := time.Parse("2006-01-02", checkOutStr)
	if err != nil {
		fmt.Println("Invalid date format. Use YYYY-MM-DD")
		return
	}

	// validate dates
	if checkOut.Before(checkIn) || checkOut.Equal(checkIn) {
		fmt.Println("Check-out date must be after check-in date ")
		return
	}

	//calculate stay duration and total cost
	nights := int(checkOut.Sub(checkIn).Hours() / 24)
	totalCost := float64(nights) * h.rooms[roomIndex].Price

	// create reservation
	reservation := Reservation{
		ID:         h.nextResID,
		RoomNumber: roomNumber,
		GuestName:  guestName,
		GuestEmail: guestEmail,
		CheckIn:    checkIn,
		CheckOut:   checkOut,
		CreatedAt:  time.Now(),
	}

	h.rooms[roomIndex].IsReserved = true
	h.reservations = append(h.reservations, reservation)
	h.nextResID++

	println("\n  Reservation Successful!")
	println("-------------------------------------")
	fmt.Printf("Rerservation ID : #%d\n", reservation.ID)
	fmt.Printf("Guest : %s (%s)\n", guestName, guestEmail)
	fmt.Printf("Room : %d (%s)\n", roomNumber, h.rooms[roomIndex].Type)
	fmt.Printf("Check-in : %s\n", checkIn.Format("2006-01-02"))
	fmt.Printf("Check-out : %s\n", checkOut.Format("2006-01-02"))
	fmt.Printf("Duration : %d night(s)\n", nights)
	fmt.Printf("Total Cost : $%.2f\n", totalCost)
	println("-------------------------------------")

}

// func for listing reservations
func (h *HotelSystem) viewReservations() {
	if len(h.reservations) == 0 {
		fmt.Println("\n No reservations found.")
		return
	}

	fmt.Println("\n Current Reservations:")
	fmt.Println("─────────────────────────────────────────────────────────────────")
	for _, res := range h.reservations {
		roomType := h.getRoomType(res.RoomNumber)
		nights := int(res.CheckOut.Sub(res.CheckIn).Hours() / 24)
		fmt.Printf("ID: #%-3d | Room: %-3d (%-8s) | Guest: %-20s\n",
			res.ID, res.RoomNumber, roomType, res.GuestName)
		fmt.Printf("         Email: %-20s | Stay: %s to %s (%d nights)\n",
			res.GuestEmail,
			res.CheckIn.Format("2006-01-02"),
			res.CheckOut.Format("2006-01-02"),
			nights)
		fmt.Println("─────────────────────────────────────────────────────────────────")
	}
}

// cancel reservation
func (h *HotelSystem) cancelReservations(reader *bufio.Reader) {
	if len(h.reservations) == 0 {
		fmt.Println("\n No reservation to cancel.")
	}

	h.viewReservations()

	resIDStr := readInput(reader, "\nEnter reservation ID to cancel: ")
	resID, err := strconv.Atoi(resIDStr)
	if err != nil {
		fmt.Println("Invalid reservation ID.")
		return
	}

	for i, res := range h.reservations {
		if resID == res.ID {
			// free up the room
			roomIndex := h.findRoomIndex(res.RoomNumber)
			if roomIndex != -1 {
				h.rooms[roomIndex].IsReserved = false
			}

			// remove the room
			h.reservations = append(h.reservations[:i], h.reservations[i+1:]...)
			fmt.Printf("Reservation #%d cancelled successfully .\n", resID)
			return
		}
	}

	fmt.Println("Reservation ID not found")
}

// helper function for room index in the list
func (h *HotelSystem) findRoomIndex(roomNumber int) int {
	for i, room := range h.rooms {
		if room.Number == roomNumber {
			return i
		}
	}

	return -1
}

// using that romIndex to get room
func (h *HotelSystem) getRoomType(roomNumber int) string {
	idx := h.findRoomIndex(roomNumber)
	if idx != -1 {
		return h.rooms[idx].Type
	}

	return "Unknown"
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (h *HotelSystem) searchRoom(reader *bufio.Reader) {
	fmt.Println("\nAvailable Room Types: Single, Double, Suite")
	roomType := readInput(reader, "Enter room type to search: ")
	roomType = cases.Title(language.English).String(strings.ToLower(roomType))

	fmt.Printf("\nAvailable %s Rooms:\n", roomType)
	fmt.Println("─────────────────────────────────────────")

	found := false
	for _, room := range h.rooms {
		if room.Type == roomType && !room.IsReserved {
			fmt.Printf("Room %d | Price: $%.2f/night\n", room.Number, room.Price)
			found = true
		}
	}

	if !found {
		fmt.Printf("No available %s rooms found.\n", roomType)
	}
	fmt.Println("─────────────────────────────────────────")
}
