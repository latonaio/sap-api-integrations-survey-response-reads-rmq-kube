package sap_api_caller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	sap_api_output_formatter "sap-api-integrations-survey-response-reads-rmq-kube/SAP_API_Output_Formatter"
	"strings"
	"sync"

	"github.com/latonaio/golang-logging-library-for-sap/logger"
	"golang.org/x/xerrors"
)

type RMQOutputter interface {
	Send(sendQueue string, payload map[string]interface{}) error
}

type SAPAPICaller struct {
	baseURL      string
	apiKey       string
	outputQueues []string
	outputter    RMQOutputter
	log          *logger.Logger
}

func NewSAPAPICaller(baseUrl string, outputQueueTo []string, outputter RMQOutputter, l *logger.Logger) *SAPAPICaller {
	return &SAPAPICaller{
		baseURL:      baseUrl,
		apiKey:       GetApiKey(),
		outputQueues: outputQueueTo,
		outputter:    outputter,
		log:          l,
	}
}

func (c *SAPAPICaller) AsyncGetSurveyResponse(iD, processorID, productID, questionUUID string, accepter []string) {
	wg := &sync.WaitGroup{}
	wg.Add(len(accepter))
	for _, fn := range accepter {
		switch fn {
		case "SurveyResponse":
			func() {
				c.SurveyResponse(iD)
				wg.Done()
			}()
		case "SurveyValuation":
			func() {
				c.SurveyValuation(processorID)
				wg.Done()
			}()
		case "SurveyValuationItem":
			func() {
				c.SurveyValuationItem(productID)
				wg.Done()
			}()
		case "SurveyQuestionAnswers":
			func() {
				c.SurveyQuestionAnswers(questionUUID)
				wg.Done()
			}()
		default:
			wg.Done()
		}
	}

	wg.Wait()
}

func (c *SAPAPICaller) SurveyResponse(iD string) {
	surveyResponseData, err := c.callSurveyResponseSrvAPIRequirementSurveyResponse("SurveyResponseRootCollectionData", iD)
	if err != nil {
		c.log.Error(err)
		return
	}
	err = c.outputter.Send(c.outputQueues[0], map[string]interface{}{"message": surveyResponseData, "function": "SurveyResponseData"})
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(surveyResponseData)

}

func (c *SAPAPICaller) callSurveyResponseSrvAPIRequirementSurveyResponse(api, iD string) ([]sap_api_output_formatter.SurveyResponse, error) {
	url := strings.Join([]string{c.baseURL, "c4codataapi", api}, "/")
	req, _ := http.NewRequest("GET", url, nil)

	c.setHeaderAPIKeyAccept(req)
	c.getQueryWithSurveyResponse(req, iD)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToSurveyResponse(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) SurveyValuation(processorID string) {
	surveyValuationData, err := c.callSurveyResponseSrvAPIRequirementSurveyValuation("SurveyResponseCollectionData", processorID)
	if err != nil {
		c.log.Error(err)
		return
	}
	err = c.outputter.Send(c.outputQueues[0], map[string]interface{}{"message": surveyValuationData, "function": "SurveyValuationData"})
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(surveyValuationData)

}

func (c *SAPAPICaller) callSurveyResponseSrvAPIRequirementSurveyValuation(api, processorID string) ([]sap_api_output_formatter.SurveyValuation, error) {
	url := strings.Join([]string{c.baseURL, "c4codataapi", api}, "/")
	req, _ := http.NewRequest("GET", url, nil)

	c.setHeaderAPIKeyAccept(req)
	c.getQueryWithSurveyValuation(req, processorID)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToSurveyValuation(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) SurveyValuationItem(productID string) {
	surveyValuationItemData, err := c.callSurveyResponseSrvAPIRequirementSurveyValuationItem("SurveyResponseItemCollectionData", productID)
	if err != nil {
		c.log.Error(err)
		return
	}
	err = c.outputter.Send(c.outputQueues[0], map[string]interface{}{"message": surveyValuationItemData, "function": "SurveyValuationItemData"})
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(surveyValuationItemData)

}

func (c *SAPAPICaller) callSurveyResponseSrvAPIRequirementSurveyValuationItem(api, productID string) ([]sap_api_output_formatter.SurveyValuationItem, error) {
	url := strings.Join([]string{c.baseURL, "c4codataapi", api}, "/")
	req, _ := http.NewRequest("GET", url, nil)

	c.setHeaderAPIKeyAccept(req)
	c.getQueryWithSurveyValuationItem(req, productID)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToSurveyValuationItem(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) SurveyQuestionAnswers(questionUUID string) {
	surveyQuestionAnswersData, err := c.callSurveyResponseSrvAPIRequirementSurveyQuestionAnswers("SurveyQuestionAnswersCollectionData", questionUUID)
	if err != nil {
		c.log.Error(err)
		return
	}
	err = c.outputter.Send(c.outputQueues[0], map[string]interface{}{"message": surveyQuestionAnswersData, "function": "SurveyQuestionAnswersData"})
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(surveyQuestionAnswersData)

}

func (c *SAPAPICaller) callSurveyResponseSrvAPIRequirementSurveyQuestionAnswers(api, questionUUID string) ([]sap_api_output_formatter.SurveyQuestionAnswers, error) {
	url := strings.Join([]string{c.baseURL, "c4codataapi", api}, "/")
	req, _ := http.NewRequest("GET", url, nil)

	c.setHeaderAPIKeyAccept(req)
	c.getQueryWithSurveyQuestionAnswers(req, questionUUID)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToSurveyQuestionAnswers(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) setHeaderAPIKeyAccept(req *http.Request) {
	req.Header.Set("APIKey", c.apiKey)
	req.Header.Set("Accept", "application/json")
}

func (c *SAPAPICaller) getQueryWithSurveyResponse(req *http.Request, iD string) {
	params := req.URL.Query()
	params.Add("$filter", fmt.Sprintf("ID eq '%s'", iD))
	req.URL.RawQuery = params.Encode()
}

func (c *SAPAPICaller) getQueryWithSurveyValuation(req *http.Request, processorID string) {
	params := req.URL.Query()
	params.Add("$filter", fmt.Sprintf("ProcessorID eq '%s'", processorID))
	req.URL.RawQuery = params.Encode()
}

func (c *SAPAPICaller) getQueryWithSurveyValuationItem(req *http.Request, productID string) {
	params := req.URL.Query()
	params.Add("$filter", fmt.Sprintf("ProductID eq '%s'", productID))
	req.URL.RawQuery = params.Encode()
}

func (c *SAPAPICaller) getQueryWithSurveyQuestionAnswers(req *http.Request, questionUUID string) {
	params := req.URL.Query()
	params.Add("$filter", fmt.Sprintf("QuestionUUID eq '%s'", questionUUID))
	req.URL.RawQuery = params.Encode()
}
