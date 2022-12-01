package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"text/template"

	_ "embed"
)

//go:embed standalone.xml.template
var configTemplateContents string

type WildflyConfiguration struct {
	ClientPort       *int              `json:"port_client"`
	DBJNDIName       *string           `json:"db_jndi_name"`
	DBMaxPoolSize    *int              `json:"db_max_pool_size"`
	DBPoolName       string            `json:"db_pool_name"`
	Metrics          bool              `json:"metrics"`
	Statistics       bool              `json:"statistics"`
	SystemProperties map[string]string `json:"systemproperties"`
}

type DatabaseConfiguration struct {
	Dialect *map[string]string `json:"dialect"`
	Name    *string            `json:"database"`
	Port    *int               `json:"port"`
	User    *string            `json:"user"`
}

type Configuration struct {
	Database DatabaseConfiguration `json:"database"`
	Wildfly  WildflyConfiguration  `json:"wildfly"`
}

type SystemProperty struct {
	Name  string
	Value string
}

type TemplateParameters struct {
	ClientPort                int
	DatabaseDriver            string
	DatabaseInitialPoolSize   int
	DatabaseMaxPoolSize       int
	DatabaseMinPoolSize       int
	DatabaseName              string
	DatabasePort              int
	DatabaseStatisticsEnabled bool
	DatabaseUser              string
	DBJNDIName                string
	SystemProperties          []SystemProperty
}

func main() {
	content, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal("Error reading stdin: ", err)
	}

	config := Configuration{}
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&config)

	//err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	tmpl, err := template.New("standalone.xml.template").Parse(configTemplateContents)
	if err != nil {
		panic(err)
	}
	var outputBuffer bytes.Buffer
	err = tmpl.Execute(&outputBuffer, getTemplateParameters(config))
	if err != nil {
		panic(err)
	}
	if !isValid(outputBuffer.String()) {
		log.Fatalf("error, invalid xml generated:\n\n%s", outputBuffer.String())
	} else {
		log.Println("Successfully generated xml")
	}
	if n, err := os.Stdout.Write(outputBuffer.Bytes()); err != nil || n < outputBuffer.Len() {
		log.Fatalf("error writing to stdout, %d bytes of %d written: %s", n, outputBuffer.Len(), err)
	}
}

func orDefault[V string | int](val *V, def V) V {
	if val != nil {
		return *val
	} else {
		return def
	}
}

func getDatabaseDriver(config Configuration) (string, string, int) {
	if config.Database.Dialect == nil {
		return "postgresql", "org.hibernate.dialect.PostgreSQL95Dialect", 5432
	}
	if len(*config.Database.Dialect) != 1 {
		log.Fatalf("expected exactly one database dialect in config, found %d", len(*config.Database.Dialect))
	}
	for key, value := range *config.Database.Dialect {
		if key == "Psql" {
			if value != "postgresql" {
				log.Fatalf("expected database dialect value to be postgresql, got %s\n", value)
			}
			return value, "org.hibernate.dialect.PostgreSQL95Dialect", 5432
		} else if key == "Mssql" {
			if value != "mssqlserver" {
				log.Fatalf("expected database dialect value to be mssqlserver, got %s\n", value)
			}
			return value, "org.hibernate.dialect.SQLServer2008Dialect", 1433
		} else {
			log.Fatalf("expected Psql or Mssql database dialect, got %s\n", key)
		}
	}
	panic("unreachable")
}

func getTemplateParameters(config Configuration) TemplateParameters {
	databaseDriver, hibernateDialect, defaultDatabasePort := getDatabaseDriver(config)
	systemProperties := []SystemProperty{
		{Name: "hibernate.dialect", Value: hibernateDialect},
	}
	for name, value := range config.Wildfly.SystemProperties {
		systemProperties = append(systemProperties, SystemProperty{Name: name, Value: value})
	}

	sort.SliceStable(systemProperties, func(i, j int) bool {
		return systemProperties[i].Name < systemProperties[j].Name
	})

	maxPoolSize := orDefault(config.Wildfly.DBMaxPoolSize, 32)

	// TODO metrics - unpick what on god's green earth this is doing and replicate:
	// https://github.com/axetrading/axetrader-installer/blob/35b8e054ce2dc855edaa9ae7d4ea72d8f4b5cf90/src/install/wildfly.rs#L914-L965

	return TemplateParameters{
		ClientPort:                orDefault(config.Wildfly.ClientPort, 8787),
		DatabaseDriver:            databaseDriver,
		DatabaseInitialPoolSize:   maxPoolSize / 8,
		DatabaseMaxPoolSize:       maxPoolSize,
		DatabaseMinPoolSize:       maxPoolSize / 8,
		DatabaseName:              orDefault(config.Database.Name, "axetrader"),
		DatabasePort:              orDefault(config.Database.Port, defaultDatabasePort),
		DatabaseStatisticsEnabled: config.Wildfly.Statistics,
		DatabaseUser:              orDefault(config.Database.User, "axetrader"),
		DBJNDIName:                orDefault(config.Wildfly.DBJNDIName, "jboss/datasources/axeDS"),
		SystemProperties:          systemProperties,
	}
}

func isValid(s string) bool {
	return xml.Unmarshal([]byte(s), new(interface{})) == nil
}
