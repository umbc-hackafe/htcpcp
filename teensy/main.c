const int ERR_UNKNOWN = -1;
const int ERR_NYI = -2;

const int ERR_KCUP_UNKNOWN = -100;
const int ERR_KCUP_EMPTY = -101;
const int ERR_KCUP_TRAY_FULL = -102;
const int ERR_KCUP_TRAY_EMPTY = -103;

const int ERR_COFFEE_DRY = -201;
const int ERR_KETTLE_DRY = -202;

const int ERR_ALREADY_BREWING = -301;

enum Addin {
  SUGAR,
  MILK,
  HONEY,
}

// Tell the turntable to seek to a particular mug
// Error if !(0 <= index <= 3)
int mug_seek(int index) {
  return ERR_NYI;
}

// Gets the current selected mug index, or the one that is currently being seeked
// Never returns an error
int mug_get() {
  return ERR_NYI;
}

// Ensure the kcup tray is open
// Repeated calls will do nothing
// Never returns an error, ideally
int kcup_tray_open() {
  return ERR_NYI;
}

// Ensure the kcup tray is closed
// Repeated calls will do nothing
// Never returns an error, ideally
int kcup_tray_close() {
  return ERR_NYI;
}

// Loads a kcup from the dispenser
// ERR_KCUP_TRAY_FULL if there is already a kcup loaded
// Opens the tray if it is closed
int kcup_load() {
  kcup_tray_open();
  return ERR_NYI;
}

// Ejects a kcup from the dispenser
// Opens the tray if it is closed
// Never returns an error
int kcup_eject() {
  kcup_tray_open();
  return ERR_NYI;
}

// Returns the number of kcups loaded
// If exact numbers are not available, returns 1 if kcups are
// available and 0 if no kcups are available
int kcup_count() {
  return ERR_NYI;
}

// Begins a brew cycle
// ERR_KCUP_TRAY_EMPTY if no kcup is loaded
// ERR_COFFEE_DRY if water reservoir is empty
// ERR_ALREADY_BREWING if a brew is already in progress
int brew() {
  return ERR_NYI;
}

// Dispenses an addin to the current mug
// Will probably not return errors
int addin_insert(enum Addin type, double amount) {
  return ERR_NYI;
}

int kettle_set_temp(int temp) {

}

int kettle_dispense(double amount) {

}
