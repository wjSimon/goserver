String inString = "";
#define D5 5
#define D6 6
void setup() {
  // put your setup code here, to run once:
  Serial.begin(115200);
  //Serial.setTimeout(20);
  pinMode(D5, OUTPUT);
  pinMode(D6, OUTPUT);
  pinMode(A0, INPUT);

  analogWrite(D5, 0);   
  analogWrite(D6, 0);   

  while (!Serial) {}
  Serial.print("test\n");
}

void loop() {
  // put your main code here, to run repeatedly:
if(Serial.available() <= 0) {return;}
  // 

  while (Serial.available() > 0) {
    int inChar = Serial.read();
    if (isDigit(inChar)) {
      // convert the incoming byte to a char
      // and add it to the string:
      inString += (char)inChar;
    }
    // if you get a newline, print the string,
    // then the string's value:
    if (inChar == '\n' || inChar == '$') {
      int pwm = inString.toInt();
      analogWrite(D5, pwm);   
      analogWrite(D6, pwm);      
      // clear the string for new input:
      inString = "";
      Serial.print(pwm); 
    }
  }
  
  //delay(1000);
}
