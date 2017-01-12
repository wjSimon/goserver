String inString = "";

void setup() {
  // put your setup code here, to run once:
  Serial.begin(115200);
  //Serial.setTimeout(20);
  pinMode(D5, OUTPUT);
  pinMode(A0, INPUT);

  while (!Serial) {}
}

void loop() {
  // put your main code here, to run repeatedly:
if(Serial.available() <= 0) {return;}
  //Serial.print("test\n"); 

  while (Serial.available() > 0) {
    int inChar = Serial.read();
    if (isDigit(inChar)) {
      // convert the incoming byte to a char
      // and add it to the string:
      inString += (char)inChar;
    }
    // if you get a newline, print the string,
    // then the string's value:
    if (inChar == '\n') {
      int pwm = inString.toInt();
      analogWrite(D5, pwm);      
      // clear the string for new input:
      inString = "";
    }
  }
  
  //delay(1000);
}
