package server

func parseWebhookRequest(h http.Header) *webhook.RequestHeaders {
	return &webhook.RequestHeaders{
		h.Get("x-github-event"),
		h.Get("x-github-delivery"),
		h.Get("x-github-signature"),
	}
}

func (h *webhook.RequestHeaders) String() string { return fmt.Sprintf("webhook event\nname: %s\nguid: %s\nsignature: %s\n", h.name, h.guid, h.signature) }

func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	headers := parseWebhookRequest(r.Header)
	
	log.Printf("handle incoming webhook:\n%+v", headers.String())

	switch headers.EventName {
	case "push":
		body, err := ioutil.ReadAll(r.Body)
	
		if err != nil {
			log.Println(err)
		}
	
		var payload *webhook.PushEventPayload
	
		err = json.Unmarshal([]byte(body), &payload)
	
		if err != nil {
			log.Println(err)
		} else {
			ci.ProcessPushEvent(*payload)
		}
		break
	default:
		break
		
	}
}