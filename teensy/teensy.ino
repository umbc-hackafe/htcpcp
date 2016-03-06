#define BUFLEN 16

#define PIN_TURNTABLE_SERVO 3
#define PIN_KCUP_LOADER 2
#define PIN_KCUP_EJECTOR 11
#define PIN_TRAY_PISTON 12
#define PIN_BREW_BUTTON 6

const int ERR_UNKNOWN = -1;
const int ERR_NYI = -2;

const int ERR_KCUP_UNKNOWN = -10;
const int ERR_KCUP_EMPTY = -11;
const int ERR_KCUP_TRAY_FULL = -12;
const int ERR_KCUP_TRAY_EMPTY = -13;

const int ERR_COFFEE_DRY = -21;
const int ERR_KETTLE_DRY = -22;

const int ERR_ALREADY_BREWING = -31;

const int ERR_MUG_INVALID = -41;

enum Addin {
  SUGAR = 0,
  MILK,
  HONEY,
}

char serial_buffer[BUFLEN] = {};
long serial_index = 0;

#define CMD_HEADER_A 0xCA
#define CMD_HEADER_B 0xFE

#define CMD_mug_seek 'a'
#define CMD_mug_get 'b'
#define CMD_kcup_tray_open 'c'
#define CMD_kcup_tray_close 'd'
#define CMD_kcup_load 'e'
#define CMD_kcup_eject 'f'
#define CMD_kcup_count 'g'
#define CMD_brew 'h'
#define CMD_addin_insert 'i'

int arglen[256] = {};
arglen[CMD_mug_seek] = 1;
arglen[CMD_mug_get] = 0;
arglen[CMD_kcup_tray_open] = 0;
arglen[CMD_kcup_tray_close] = 0;
arglen[CMD_kcup_load] = 0;
arglen[CMD_kcup_eject] = 0;
arglen[CMD_kcup_count] = 0;
arglen[CMD_brew] = 0;
arglen[CMD_addin_insert] = 2;

#define TRAY_OPEN 0
#define TRAY_CLOSED 1

#define LOADER_REST 0
#define LOADER_ACTIVE 1

#define EJECTOR_OFF 0
#define EJECTOR_ON 1

int mug_index = 0;
int tray_state = 2;
int ejector_state = 2;
int kcup_loaded = 1;

void send_data(int result) {
  Serial.write(result);
}

void handle_serial {
  if (Serial.available()) {
    int data = Serial.read();
    if (data == CMD_HEADER_A) {
      serial_index = 0;
    } else {
      serial_buffer[serial_index++] = data;
    }
  }

  if (serial_buffer[0] == CMD_HEADER_B && serial_index > 1) {
    if (serial_index > 1 + arglen[serial_buffer[1]]) {
      int res;
      switch (serial_buffer[1]) {
      case CMD_mug_seek:
	// int
	res = mug_seek(serial_buffer[2]);
	break;
      case CMD_mug_get:
	res = mug_get();
	break;
      case CMD_kcup_tray_open:
	res = kcup_tray_open();
	break;
      case CMD_kcup_tray_close:
	res = kcup_tray_close();
	break;
      case CMD_kcup_load:
	res = kcup_load();
	break;
      case CMD_kcup_eject:
	res = kcup_eject();
	break;
      case CMD_kcup_count:
	res = kcup_count();
	break;
      case CMD_brew:
	res = brew();
	break;
      case CMD_addin_insert:
	// int, int
	res = brew(serial_buffer[2], serial_buffer[3]);
	break;
      default:
	res = ERR_UNKNOWN;
	break;
      }
      send_data(res);
      serial_index = 0;
    }
  }
}

#define MUG_ANGLE_MIN 0
#define MUG_ANGLE_MAX 255
#define MUG_INDEX_MIN 0
#define MUG_INDEX_MAX 5

// Tell the turntable to seek to a particular mug
// Error if !(0 <= index <= 5)
int mug_seek(int index) {
  if (index < MUG_INDEX_MIN || index > MUG_INDEX_MAX)
    return ERR_MUG_INVALID;

  int angle = (MUG_ANGLE_MAX - MUG_ANGLE_MIN) / (MUG_INDEX_MAX - MUG_INDEX_MIN) * (index - MUG_INDEX_MIN) + MUG_ANGLE_MIN;
  analogWrite(PIN_TURNTABLE_SERVO, angle);
  mug_index = index;
  return 0;
}

// Gets the current selected mug index, or the one that is currently being seeked
// Never returns an error
int mug_get() {
  return mug_index;
}

// Ensure the kcup tray is open
// Repeated calls will do nothing
// Never returns an error, ideally
int kcup_tray_open() {
  digitalWriteFast(PIN_TRAY_PISTON, TRAY_OPEN);
  if (tray_state != TRAY_OPEN) {
    delay(500);
  }
  tray_state = TRAY_OPEN;
  KCUP_LOADER, KCUP_EJECTOR, TRAY_PISTON
  return 0;
}

// Ensure the kcup tray is closed
// Repeated calls will do nothing
// Never returns an error, ideally
int kcup_tray_close() {
  digtalWriteFast(PIN_TRAY_PISTON, TRAY_CLOSED);
  if (tray_state != TRAY_CLOSED) {
    delay(500);
  }
  tray_state = TRAY_CLOSED;
  return 0;
}

// Loads a kcup from the dispenser
// ERR_KCUP_TRAY_FULL if there is already a kcup loaded
// Opens the tray if it is closed
int kcup_load() {
  if (kcup_loaded) {
    return ERR_KCUP_TRAY_FULL;
  }

  kcup_tray_open();

  digitalWriteFast(PIN_KCUP_LOADER, LOADER_ACTIVE);
  delay(500);
  digitalWriteFast(PIN_KCUP_LOADER, LOADER_REST);
  delay(500);

  kcup_loaded = 1;

  return 0;
}

// Ejects a kcup from the dispenser
// Opens the tray if it is closed
// Never returns an error
int kcup_eject() {
  kcup_tray_open();

  kcup_loaded = 0;

  digitalWriteFast(PIN_KCUP_EJECTOR, EJECTOR_ON);
  delay(100);
  digitalWriteFast(PIN_KCUP_EJECTOR, EJECTOR_OFF);
  delay(250);

  return 0;
}

// Returns the number of kcups loaded
// If exact numbers are not available, returns 1 if kcups are
// available and 0 if no kcups are available
int kcup_count() {
  return 1;
}

// Begins a brew cycle
// ERR_COFFEE_DRY if water reservoir is empty (not yet)
// ERR_ALREADY_BREWING if a brew is already in progress (not yet)
int brew() {
  if (!kcup_loaded) kcup_load();

  digitalWriteFast(PIN_BREW_BUTTON, 1);
  delay(100);
  digitalWriteFast(PIN_BREW_BUTTON, 0);

  return 0;
}

// Dispenses an addin to the current mug
// Will probably not return errors
int addin_insert(enum Addin type, int amount) {
  return ERR_NYI;
}

int kettle_set_temp(int temp) {

}

int kettle_dispense(int amount) {

}

void initialize() {
  analogWriteResolution(8);
  pinMode(PIN_TURNTABLE_SERVO, OUTPUT);
  pinMode(PIN_KCUP_LOADER, OUTPUT);
  pinMode(PIN_KCUP_EJECTOR, OUTPUT);
  pinMode(PIN_TRAY_PISTON, OUTPUT);

  Serial.begin(9600);
}

int main() {
  initialize();

  while (1) {
    handle_serial();
  }
}
