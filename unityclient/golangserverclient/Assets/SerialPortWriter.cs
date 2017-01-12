using UnityEngine;
using System.Collections;
using System.IO.Ports;

public class SerialPortWriter : MonoBehaviour {

    SerialPort sp = new SerialPort("COM4", 115200);

    int counter = 0;
	// Use this for initialization
	void Start () {
        sp.Open();

	}
	
	// Update is called once per frame
	void Update () {
        //Debug.Log("sp open? " + sp.IsOpen);

        counter+=10;
        counter %= 1024;
        try
        {
            sp.Write("" + counter + "\n");
        } catch (System.Exception e)
        {
            if(sp.IsOpen) sp.Close();
            Debug.Log("reconnect...");
            try
            {
                sp.Open();
            } catch (System.Exception e2)
            {
                //
            }
        }

    }
}
