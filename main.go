package main

import (
	"context"
	"fmt"
	"log"
	"os"

	keto "github.com/ory/keto-client-go"
)

func CreateRelationTuple(writeCLI *keto.APIClient, subject, relation, namespace, object string) error {
	relationQuery := keto.NewRelationQuery()
	relationQuery.Namespace = &namespace
	relationQuery.Object = &object
	relationQuery.SubjectId = &subject
	relationQuery.Object = &object
	relationQuery.Relation = &relation

	_, _, err := writeCLI.WriteApi.CreateRelationTuple(context.Background()).RelationQuery(*relationQuery).Execute()

	if err != nil {
		return fmt.Errorf("error while creating relation tuple: %w", err)
	}

	return nil
}

func CreateRelationTupleWithSubjectSet(writeCLI *keto.APIClient, subjectNamespace, subjectObject, subjectRelation, relation, namespace, object string) error {
	relationQuery := keto.NewRelationQuery()
	subjectSet := keto.NewSubjectSet(subjectNamespace, subjectObject, subjectRelation)
	relationQuery.Namespace = &namespace
	relationQuery.Object = &object
	relationQuery.SubjectSet = subjectSet
	relationQuery.Object = &object
	relationQuery.Relation = &relation

	_, _, err := writeCLI.WriteApi.CreateRelationTuple(context.Background()).RelationQuery(*relationQuery).Execute()

	if err != nil {
		return fmt.Errorf("error while creating relation tuple: %w", err)
	}

	return nil
}

func CheckRelationTuple(readCLI *keto.APIClient, namespace, object, relation, subjectId string) (bool, error) {
	resp, httpResp, err := readCLI.ReadApi.GetCheck(context.Background()).Namespace(namespace).Object(object).Relation(relation).SubjectId(subjectId).Execute()

	if err != nil {
		if httpResp.StatusCode == 403 {
			return false, nil
		}
		return false, fmt.Errorf("error while checking relation tuple: %w", err)
	}

	return resp.Allowed, nil
}

func main() {
	writeConfig := keto.NewConfiguration()
	writeConfig.Host = "127.0.0.1:4467"
	writeConfig.Scheme = "http"
	writeCLI := keto.NewAPIClient(writeConfig)

	readConfig := keto.NewConfiguration()
	readConfig.Host = "127.0.0.1:4466"
	readConfig.Scheme = "http"
	readCLI := keto.NewAPIClient(readConfig)

	// Make Bob the owner of notebook1
	if err := CreateRelationTuple(writeCLI, "bob", "owner", "notebooks", "notebook1"); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// Bob also owns codeinsight1
	if err := CreateRelationTuple(writeCLI, "bob", "owner", "codeinsights", "codeinsight1"); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// The first three params form a subject set of all the owners of notebook1.
	// So now we're saying "Everyone that is an owner of notebook1 is also an owner of codeinsight1"
	// Essentially we're coupling these permissions and codeinsight1 is "embedded" in notebook1
	if err := CreateRelationTupleWithSubjectSet(writeCLI, "notebooks", "notebook1", "owner", "owner", "codeinsights", "codeinsight1"); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// Next we want to show how we can use this coupling to share groups of resources to groups of users

	// We share the notebook with everyone in dragonteam
	// "Anyone that is a member of dragonteam is an owner of notebook1"
	if err := CreateRelationTupleWithSubjectSet(writeCLI, "teams", "dragonteam", "member", "owner", "notebooks", "notebook1"); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// User Steven is now a member of dragonteam
	if err := CreateRelationTuple(writeCLI, "steven", "member", "teams", "dragonteam"); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// As a result, because codeinsight1 is embedded in notebook1, and notebook1 has been shared with dragonteam
	// and Steven is a member of dragonteam, Steven should be a valid owner of codeinsight1

	canRead, err := CheckRelationTuple(readCLI, "codeinsights", "codeinsight1", "owner", "steven")

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if canRead {
		log.Println("Steven can read Code Insight 1")
	} else {
		log.Println("Steven cannot read Code Insight 1")
	}

	// Running this code should have printed "Steven can read Code Insight 1"
}
