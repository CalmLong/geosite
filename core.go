package main

import (
	"errors"
	"strconv"
	"strings"
	
	"github.com/v2fly/v2ray-core/v4/app/router"
)

type Entry struct {
	Type  string
	Value string
	Attrs []*router.Domain_Attribute
}

type List struct {
	Name  string
	Entry []Entry
}

type ParsedList struct {
	Name      string
	Inclusion map[string]bool
	Entry     []Entry
}

func (l *ParsedList) toProto() (*router.GeoSite, error) {
	site := &router.GeoSite{
		CountryCode: l.Name,
	}
	for _, entry := range l.Entry {
		switch entry.Type {
		case "domain":
			site.Domain = append(site.Domain, &router.Domain{
				Type:      router.Domain_Domain,
				Value:     entry.Value,
				Attribute: entry.Attrs,
			})
		case "regexp":
			site.Domain = append(site.Domain, &router.Domain{
				Type:      router.Domain_Regex,
				Value:     entry.Value,
				Attribute: entry.Attrs,
			})
		case "keyword":
			site.Domain = append(site.Domain, &router.Domain{
				Type:      router.Domain_Plain,
				Value:     entry.Value,
				Attribute: entry.Attrs,
			})
		case "full":
			site.Domain = append(site.Domain, &router.Domain{
				Type:      router.Domain_Full,
				Value:     entry.Value,
				Attribute: entry.Attrs,
			})
		default:
			return nil, errors.New("unknown domain type: " + entry.Type)
		}
	}
	return site, nil
}

func parseDomain(domain string, entry *Entry) error {
	kv := strings.Split(domain, ":")
	if len(kv) == 1 {
		entry.Type = "domain"
		entry.Value = strings.ToLower(kv[0])
		return nil
	}
	
	if len(kv) == 2 {
		entry.Type = strings.ToLower(kv[0])
		entry.Value = strings.ToLower(kv[1])
		return nil
	}
	
	return errors.New("Invalid format: " + domain)
}

func parseAttribute(attr string) (*router.Domain_Attribute, error) {
	var attribute router.Domain_Attribute
	if len(attr) == 0 || attr[0] != '@' {
		return &attribute, errors.New("invalid attribute: " + attr)
	}
	
	// Trim attribute prefix `@` character
	attr = attr[1:]
	parts := strings.Split(attr, "=")
	if len(parts) == 1 {
		attribute.Key = strings.ToLower(parts[0])
		attribute.TypedValue = &router.Domain_Attribute_BoolValue{BoolValue: true}
	} else {
		attribute.Key = strings.ToLower(parts[0])
		intv, err := strconv.Atoi(parts[1])
		if err != nil {
			return &attribute, errors.New("invalid attribute: " + attr + ": " + err.Error())
		}
		attribute.TypedValue = &router.Domain_Attribute_IntValue{IntValue: int64(intv)}
	}
	return &attribute, nil
}

func parseEntry(line string) (Entry, error) {
	line = strings.TrimSpace(line)
	parts := strings.Split(line, " ")
	
	var entry Entry
	if len(parts) == 0 {
		return entry, errors.New("empty entry")
	}
	
	if err := parseDomain(parts[0], &entry); err != nil {
		return entry, err
	}
	
	for i := 1; i < len(parts); i++ {
		attr, err := parseAttribute(parts[i])
		if err != nil {
			return entry, err
		}
		entry.Attrs = append(entry.Attrs, attr)
	}
	
	return entry, nil
}

func isMatchAttr(Attrs []*router.Domain_Attribute, includeKey string) bool {
	isMatch := false
	mustMatch := true
	matchName := includeKey
	if strings.HasPrefix(includeKey, "!") {
		isMatch = true
		mustMatch = false
		matchName = strings.TrimLeft(includeKey, "!")
	}
	
	for _, Attr := range Attrs {
		attrName := Attr.Key
		if mustMatch {
			if matchName == attrName {
				isMatch = true
				break
			}
		} else {
			if matchName == attrName {
				isMatch = false
				break
			}
		}
	}
	return isMatch
}

func createIncludeAttrEntrys(list *List, matchAttr *router.Domain_Attribute) []Entry {
	newEntryList := make([]Entry, 0, len(list.Entry))
	matchName := matchAttr.Key
	for _, entry := range list.Entry {
		matched := isMatchAttr(entry.Attrs, matchName)
		if matched {
			newEntryList = append(newEntryList, entry)
		}
	}
	return newEntryList
}

func ParseList(list *List, ref map[string]*List) (*ParsedList, error) {
	pl := &ParsedList{
		Name:      list.Name,
		Inclusion: make(map[string]bool),
	}
	entryList := list.Entry
	for {
		newEntryList := make([]Entry, 0, len(entryList))
		hasInclude := false
		for _, entry := range entryList {
			if entry.Type == "include" {
				refName := strings.ToUpper(entry.Value)
				if entry.Attrs != nil {
					for _, attr := range entry.Attrs {
						InclusionName := strings.ToUpper(refName + "@" + attr.Key)
						if pl.Inclusion[InclusionName] {
							continue
						}
						pl.Inclusion[InclusionName] = true
						
						refList := ref[refName]
						if refList == nil {
							return nil, errors.New(entry.Value + " not found.")
						}
						attrEntrys := createIncludeAttrEntrys(refList, attr)
						if len(attrEntrys) != 0 {
							newEntryList = append(newEntryList, attrEntrys...)
						}
					}
				} else {
					InclusionName := refName
					if pl.Inclusion[InclusionName] {
						continue
					}
					pl.Inclusion[InclusionName] = true
					refList := ref[refName]
					if refList == nil {
						return nil, errors.New(entry.Value + " not found.")
					}
					newEntryList = append(newEntryList, refList.Entry...)
				}
				hasInclude = true
			} else {
				newEntryList = append(newEntryList, entry)
			}
		}
		entryList = newEntryList
		if !hasInclude {
			break
		}
	}
	pl.Entry = entryList
	
	return pl, nil
}