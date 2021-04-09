package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	//"sync"
	//"time"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/piprate/json-gold/ld"
//	"github.com/satori/go.uuid"
//	. "github.com/smartfog/fogflow/common/config"
	. "github.com/smartfog/fogflow/common/constants"
//	. "github.com/smartfog/fogflow/common/datamodel"
	. "github.com/smartfog/fogflow/common/ngsi"
)

type LdThinBroker struct {
	SecurityCfg     *HTTPS
	myProfile BrokerProfile

        //NGSI-LD feature addition
        ldEntities      map[string]interface{} // to map Entity Id with LDContextElement.
        ldEntities_lock sync.RWMutex

        ldContextRegistrations      map[string]CSourceRegistrationRequest // to map Registration Id with CSourceRegistrationRequest.
        ldContextRegistrations_lock sync.RWMutex

        ldEntityID2RegistrationID      map[string]string //to map the Entity IDs with their registration id.
        ldEntityID2RegistrationID_lock sync.RWMutex

        ldSubscriptions      map[string]*LDSubscriptionRequest // to map Subscription Id with LDSubscriptionRequest.
        ldSubscriptions_lock sync.RWMutex

        tmpNGSIldNotifyCache    []string
        tmpNGSILDNotifyCache    map[string]*NotifyContextAvailabilityRequest
        entityId2LDSubcriptions map[string][]string

}

//NGSILD upsert API

func (ldTb *LdThinBroker) LDUpdateContext(w rest.ResponseWriter, r *rest.Request) {
	err := contentTypeValidator(r.Header.Get("Content-Type"))
	if err != nil {
		w.WriteHeader(500)
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		reqBytes, _ := ioutil.ReadAll(r.Body)
		var LDupdateCtxReq []interface{}

		err = json.Unmarshal(reqBytes, &LDupdateCtxReq)

		if err != nil {
			err := errors.New("This interface only supports arrays of entities")
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res := ResponseError{}

		for _, ctx := range LDupdateCtxReq {
			var context []interface{}
			contextInPayload := true
			//Get Link header if present
			Link := r.Header.Get("Link")
			if link := r.Header.Get("Link"); link != "" {
				contextInPayload = false                    // Context in Link header
				linkMap := ldTb.extractLinkHeaderFields(link) // Keys in returned map are: "link", "rel" and "type"
				if linkMap["rel"] != DEFAULT_CONTEXT {
				}
			}
			context = append(context, DEFAULT_CONTEXT)

			//Get a resolved object ([]interface object)
			resolved, err := ldTb.ExpandPayload(ctx, context, contextInPayload)
			if err != nil {

				if err.Error() == "EmptyPayload!" {
					//res.Errors.Details  = "EmptyPayload is not allowed"
					problemSet := ProblemDetails{}
					problemSet.Details = "EmptyPayload is not allowed"
					res.Errors = append(res.Errors, problemSet)
					continue
				}
				if err.Error() == "Id can not be nil!" {
					problemSet := ProblemDetails{}
					problemSet.Details = "Id can not be nil!"
					res.Errors = append(res.Errors, problemSet)
					continue
				}
				if err.Error() == "Type can not be nil!" {
					problemSet := ProblemDetails{}
					problemSet.Details = "Type can not be nil!"
					res.Errors = append(res.Errors, problemSet)
					continue
				}
				//res.Errors.Details  = "Unknown"
				problemSet := ProblemDetails{}
				problemSet.Details = "Unkown!"
				res.Errors = append(res.Errors, problemSet)
				continue
			} else {
				sz := Serializer{}

				// Deserialize the payload here.
				deSerializedEntity, err := sz.DeSerializeEntity(resolved)
				if err != nil {
					problemSet := ProblemDetails{}
					problemSet.Details = "Problem in deserialization"
					res.Errors = append(res.Errors, problemSet)
					continue
				} else {
					//Update createdAt value.
					if !strings.HasPrefix(deSerializedEntity["id"].(string), "urn:ngsi-ld:") {
						problemSet := ProblemDetails{}
						problemSet.Details = "Entity id must contain uri!"
						res.Errors = append(res.Errors, problemSet)
						continue
					}
					deSerializedEntity["@context"] = context
					fmt.Println(deSerializedEntity)
					res.Success = append(res.Success, deSerializedEntity["id"].(string))
					ldTb.handleLdExternalUpdateContext(deSerializedEntity, Link)
				}
			}
		}
		if res.Errors != nil && res.Success == nil {
			w.WriteHeader(404)
			w.WriteJson(&res)
		}
		if res.Errors != nil && res.Success != nil {
			w.WriteHeader(207)
			w.WriteJson(&res)
		}
		if res.Errors == nil && res.Success != nil {
			w.WriteHeader(http.StatusNoContent)
		}

	}
}

func (ldTb *LdThinBroker) UpdateLdContext2LocalSite(updateCtxReq map[string]interface{}) {
        ldTb.ldEntities_lock.Lock()
        eid := getId(updateCtxReq)
        hasLdUpdatedMetadata := hasLdUpdatedMetadata(updateCtxReq, ldTb.ldEntities[eid])
        ldTb.ldEntities_lock.Unlock()

        ldTb.updateLdContextElement(updateCtxReq)

        go ldTb.LDNotifySubscribers(updateCtxReq, true)
        if hasLdUpdatedMetadata == true {
                ldTb.registerLDContextElement(updateCtxReq)
        }
}

func (ldTb *LdThinBroker) UpdateLdContext2RemoteSite(updateCtxReq map[string]interface{}, brokerURL string, link string) {
        INFO.Println(brokerURL)
        client := NGSI10Client{IoTBrokerURL: brokerURL, SecurityCfg: ldTb.SecurityCfg}
        client.CreateLDEntityOnRemote(updateCtxReq, link)
}


func (ldTb *LdThinBroker) queryOwnerOfLDEntity(eid string) string {
        inLocalBroker := true

        ldTb.ldEntities_lock.RLock()
        _, exist := tb.ldEntities[eid]
        inLocalBroker = exist
	ldTb.ldEntities_lock.RUnlock()

        if inLocalBroker == true {
                return ldTb.myProfile.MyURL
        } else {
                client := NGSI9Client{IoTDiscoveryURL: tb.IoTDiscoveryURL, SecurityCfg: ldTb.SecurityCfg}
                brokerURL, _ := client.GetProviderURL(eid)
                if brokerURL == "" {
                        return ldTb.myProfile.MyURL
                }
                return brokerURL
        }
}

func (ldTb *LdThinBroker) handleLdExternalUpdateContext(updateCtxReq map[string]interface{}, link string) {
        eid := getId(updateCtxReq)
        brokerURL := tb.queryOwnerOfLDEntity(eid)
        if brokerURL == tb.myProfile.MyURL {
                ldTb.UpdateLdContext2LocalSite(updateCtxReq)
        } else {
                ldTb.UpdateLdContext2RemoteSite(updateCtxReq, brokerURL, link)
        }
}


// Expand the payload
func (ldTb *LdThinBroker) ExpandPayload(ctx interface{}, context []interface{}, contextInPayload bool) ([]interface{}, error) {
        //get map[string]interface{} of reqBody
        itemsMap, err := ldTb.getStringInterfaceMap(ctx)
        if err != nil {
                return nil, err
        } else {
                // Check the type of payload: Entity, registration or Subscription
                var payloadType string
                if _, ok := itemsMap["type"]; ok == true {
                        payloadType = itemsMap["type"].(string)
                } else if _, ok := itemsMap["@type"]; ok == true {
                        typ := itemsMap["@type"].([]interface{})
                        payloadType = typ[0].(string)
                }
                if payloadType == "" {
                        err := errors.New("Type can not be nil!")
                        return nil, err
                }
                if payloadType != "ContextSourceRegistration" && payloadType != "Subscription" {
                        // Payload is of Entity Type
                        // Check if some other broker is registered for providing this entity or not
                        var entityId string
                        if _, ok := itemsMap["id"]; ok == true {
                                entityId = itemsMap["id"].(string)
                        } else if _, ok := itemsMap["@id"]; ok == true {
                                entityId = itemsMap["@id"].(string)
                        }

                        if entityId == "" {
                                err := errors.New("Id can not be nil!")
                                return nil, err
                        }
                        if contextInPayload == true {
                                if Context := itemsMap["@context"]; Context == nil {
                                        err := errors.New("@context is Empty")
                                        return nil, err
                                }
                                if Context := itemsMap["@context"].([]interface{}); len(Context) == 0 {
                                        err := errors.New("@context is Empty")
                                        return nil, err
                                }
                        }
                }

                // Update Context in itemMap
                if contextInPayload == true && itemsMap["@context"] != nil {
                        contextItems := itemsMap["@context"].([]interface{})
                        context = append(context, contextItems...)
                }
                itemsMap["@context"] = context

                if expanded, err := ldTb.ExpandData(itemsMap); err != nil {
                        return nil, err
                } else {

                        return expanded, nil
                }
        }
}

// Expand the NGSI-LD Data with context
func (ldTb *LdThinBroker) ExpandData(v interface{}) ([]interface{}, error) {
        proc := ld.NewJsonLdProcessor()
        options := ld.NewJsonLdOptions("")
        //LD processor expands the data and returns []interface{}
        expanded, err := proc.Expand(v, options)
        return expanded, err
}

//Get string-interface{} map from request body
func (ldTb *LdThinBroker) getStringInterfaceMap(ctx interface{}) (map[string]interface{}, error) {
        itemsMap := ctx.(map[string]interface{})

        if len(itemsMap) != 0 {
                return itemsMap, nil
        } else {
                return nil, errors.New("EmptyPayload!")
        }
}


func (ldTb *LdThinBroker) extractLinkHeaderFields(link string) map[string]string {
        mp := make(map[string]string)
        linkArray := strings.Split(link, ";")

        for i, arrValue := range linkArray {
                linkArray[i] = strings.Trim(arrValue, " ")
                if strings.HasPrefix(arrValue, "<{{link}}>") {
                        continue // TBD, context link
                } else if strings.HasPrefix(arrValue, "http") {
                        mp["link"] = arrValue
                } else if strings.HasPrefix(arrValue, " rel=") {
                        mp["rel"] = arrValue[6 : len(arrValue)-1] // Trimmed `rel="` and `"`
                } else if strings.HasPrefix(arrValue, " type=") {
                        mp["type"] = arrValue[7 : len(arrValue)-1] // Trimmed `type="` and `"`
                }
        }

        return mp
}
