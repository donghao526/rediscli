package rediscli

import (
	"fmt"
	"strings"
)

// get the resp protocol string for the command
func GetRespStrOfCmd(command string) string {
	arrKeyWords := getKeyWordsOfCmd(command)
	return buildRespString(arrKeyWords)
}

// get key words of the cmd
func getKeyWordsOfCmd(command string) []string {
	var keywords = strings.Fields(command)
	return keywords
}

// build the resp protocol string for the command
func buildRespString(arrKeyWords []string) string {
	argc := len(arrKeyWords)
	var respStr = ""
	respStr += fmt.Sprintf("*%d\r\n", argc)

	for _, value := range arrKeyWords {
		if value != "" {
			respStr += fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
		}
	}
	return respStr
}

// parse reply
func ParseReply(reply *RedisObject) string {
	out := ""
	switch reply.obj_type {
	case TYPE_STRING:
		fallthrough
	case TYPE_ERROR:
		out += reply.str_value
	case TYPE_BULK:
		out += fmt.Sprintf("\"%s\"", reply.str_value)
	case TYPE_INTEGER:
		out += fmt.Sprintf("(integer) %d", reply.int_value)
	case TYPE_NIL:
		out += fmt.Sprintf("(integer) %s", reply.str_value)
	case TYPE_ARRAY:
		out = parseArray(reply)
	}
	return out
}

func parseArray(reply *RedisObject) string {
	out := ""
	i := 0
	format := "%d) %s\n"
	for i < reply.size {
		if i == reply.size - 1 {
			format = "%d) %s"
		}
		out += fmt.Sprintf(format, i + 1, ParseReply(reply.member[i]))
		i++
	}
	return out
}

func isSpace(char byte) bool {
	if char == ' ' {
		return true
	} else {
		return false
	}
}

func splitArgs(line string) []string {
	var args []string
	argc := 0
	p := 0
	length := len(line)
	for true {
		/* skip blanks */
		for p < length && isSpace(line[p]) {
			p++
		}

		if p < length {
			indq := false
			insq := false
			done := false

			current := ""

			for !done {
				if indq == true {
					if line[p] == '\\' && line[p+1] == 'x' {
						p += 3
					} else if line[p] == '\\' && p + 1 < length {
						p++
						switch line[p] {
						case 'n': current += string('\n')
						case 'r': current += string('\r')
						case 't': current += string('\t')
						case 'b': current += string('\b')
						case 'a': current += string('\a')
						default:
							current += string(line[p])
						}
					} else if line[p] == '"' {
						if p + 1 < length && !isSpace(line[p+1]) {
							done = true
						}
					} else if p >= length {
						goto err
					} else {
					}
				} else if insq == true {
					if line[p] == '\\' && line[p+1] == '\'' {
						p++
						current += "'"
					} else if line[p] == '\'' {
						/* closing quote must be followed by a space or
						 * nothing at all. */
						if p + 1 < length && !isSpace(line[p+1]) {
							goto err
						}
						done = true
					} else if p >= length {
						/* unterminated quotes */
						goto err
					} else {
						current += string(line[p])
					}
				} else {
					switch line[p] {
					case ' ': fallthrough
					case '\n': fallthrough
					case '\r': fallthrough
					case '\t': fallthrough
					case '0':
						done = true
					case '"':
						indq = true
					case '\'':
						insq = true
					default:
						current = string(line[p])
					}
				}
				p++
			}
		}
	}

err:
	return nil
}