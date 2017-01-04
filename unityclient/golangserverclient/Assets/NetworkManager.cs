using UnityEngine;
using System.Collections;
using System.Collections.Generic;
using System.Net;

public class NetworkManager : MonoBehaviour {

	string username = "test";
	string password = "123456";
	string host = "http://localhost";

	// Use this for initialization
	void Start () {
		StartCoroutine(Login());
	}
	
	// Update is called once per frame
	void Update () {
	
	}

	IEnumerator Login()
	{
		string url = host + "/rpc/loginunity";
		WWWForm form = new WWWForm();

		form.AddField("Name", username);
		form.AddField("Password", password);

		WWW www = new WWW(url, form);
		yield return www;

		Debug.Log(GetResponseCode(www));
		Debug.Log(GetCookie(www));
		/**/
	}


	public static string GetCookie(WWW request)
	{
		string cookie = "";
		foreach (KeyValuePair<string, string> entry in request.responseHeaders)
		{
			Debug.Log(entry.Key + "=" + entry.Value);
		}

		request.responseHeaders.TryGetValue("Set-Cookie", out cookie);
		return cookie;
	}
	public static int GetResponseCode(WWW request)
	{
		int ret = 0;
		if (request.responseHeaders == null)
		{
			Debug.LogError("no response headers.");
		}
		else
		{
			if (!request.responseHeaders.ContainsKey("STATUS"))
			{
				Debug.LogError("response headers has no STATUS.");
			}
			else
			{
				ret = parseResponseCode(request.responseHeaders["STATUS"]);
			}
		}

		return ret;
	}

	public static int parseResponseCode(string statusLine)
	{
		int ret = 0;

		string[] components = statusLine.Split(' ');
		if (components.Length < 3)
		{
			Debug.LogError("invalid response status: " + statusLine);
		}
		else
		{
			if (!int.TryParse(components[1], out ret))
			{
				Debug.LogError("invalid response code: " + components[1]);
			}
		}

		return ret;
	}
}
