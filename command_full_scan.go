package pegic

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"pegic/ast"
	p "pegic/parser"

	"github.com/XiaoMi/pegasus-go-client/pegasus"
	"github.com/olekukonko/tablewriter"
)

type fullScanCommand struct {
	count   bool
	delete  bool
	hashKey *fullScanHashKey
	sortKey *fullScanSortKey
	noValue bool
}

type fullScanHashKey struct {
	contains string
	suffix   string
	prefix   string
}

type fullScanSortKey struct {
	is       string
	contains string
	suffix   string
	prefix   string
	start    string
	stop     string
}

func (c *fullScanCommand) parse(input string) error {
	res, s := p.ArrayWhiteSpace(
		p.Opt(p.Alt(p.TagNoCase("count"), p.TagNoCase("delete"))),
		p.Opt(
			p.ArrayWhiteSpace(
				p.TagNoCase("hashkey"),
				p.Alt(p.TagNoCase("contains"), p.TagNoCase("suffix"), p.TagNoCase("prefix")),
				p.String,
			),
		),
		p.Opt(
			p.ArrayWhiteSpace(
				p.TagNoCase("sortkey"),
				p.Opt(p.Alt(
					p.ArrayWhiteSpace(p.Alt(p.TagNoCase("is"), p.TagNoCase("contains"), p.TagNoCase("suffix"), p.TagNoCase("prefix")), p.String),
					p.ArrayWhiteSpace(p.TagNoCase("between"), p.String, p.TagNoCase("and"), p.String),
				)),
			),
		),
		p.Opt(p.TagNoCase("novalue")),
	)(input)
	if res.Err != nil {
		return res.Err
	}
	if s != "" {
		return fmt.Errorf("redundant input `%s`", s)
	}
	out := res.Output.([]interface{})
	if out[0] != nil {
		c.count = out[0].(string) == "count"
		c.delete = !c.count
	} else {
		c.count = false
		c.delete = false
	}
	if out[1] != nil {
		hk := &fullScanHashKey{}
		p := out[1].([]interface{})
		switch p[1].(string) {
		case "contains":
			hk.contains = p[2].(string)
		case "suffix":
			hk.suffix = p[2].(string)
		case "prefix":
			hk.prefix = p[2].(string)
		}
		c.hashKey = hk
	} else {
		c.hashKey = nil
	}
	if out[2] != nil {
		sk := &fullScanSortKey{}
		p := out[2].([]interface{})
		if p[1] != nil {
			a := p[1].([]interface{})
			if a[0] == "between" {
				sk.start = a[1].(string)
				sk.stop = a[3].(string)
			} else {
				switch a[0].(string) {
				case "is":
					sk.is = a[1].(string)
				case "contains":
					sk.contains = a[1].(string)
				case "prefix":
					sk.prefix = a[1].(string)
				case "suffix":
					sk.suffix = a[1].(string)
				}
			}
		}
		c.sortKey = sk
	} else {
		c.sortKey = nil
	}
	c.noValue = out[3] != nil
	return nil
}

func (c *fullScanCommand) execute(ctx *ExecContext) error {
	if ctx.table == nil {
		return noTableError
	}
	hkFilter := pegasus.Filter{}
	if hk := c.hashKey; hk != nil {
		if hk.contains != "" {
			hkFilter.Type = pegasus.FilterTypeMatchAnywhere
			hkFilter.Pattern = []byte(hk.contains)
		}
		if hk.prefix != "" {
			hkFilter.Type = pegasus.FilterTypeMatchPrefix
			hkFilter.Pattern = []byte(hk.prefix)
		}
		if hk.suffix != "" {
			hkFilter.Type = pegasus.FilterTypeMatchPostfix
			hkFilter.Pattern = []byte(hk.suffix)
		}
	}
	var (
		start []byte
		stop  []byte
		is    []byte
	)
	skFilter := pegasus.Filter{}
	if sk := c.sortKey; sk != nil {
		if sk.contains != "" {
			skFilter.Type = pegasus.FilterTypeMatchAnywhere
			skFilter.Pattern = []byte(sk.contains)
		}
		if sk.prefix != "" {
			skFilter.Type = pegasus.FilterTypeMatchPrefix
			skFilter.Pattern = []byte(sk.prefix)
		}
		if sk.suffix != "" {
			skFilter.Type = pegasus.FilterTypeMatchPostfix
			skFilter.Pattern = []byte(sk.suffix)
		}
		if sk.is != "" {
			is = []byte(sk.is)
		}
		if sk.start != "" {
			start = []byte(sk.start)
			stop = []byte(sk.stop)
		}
	}
	sopts := &pegasus.ScannerOptions{
		HashKeyFilter: hkFilter,
		SortKeyFilter: skFilter,
		NoValue:       c.noValue,
	}
	scanners, err := ctx.table.GetUnorderedScanners(context.Background(), 16, sopts)
	if err != nil {
		return err
	}
	var result [][]string
	keys := make(map[string][][]byte)
	for _, scanner := range scanners {
		for {
			completed, hashKey, sortKey, value, err := scanner.Next(context.Background())
			if err != nil {
				return err
			}
			if completed {
				break
			}
			if is != nil && !bytes.Equal(is, sortKey) {
				continue
			}
			if start != nil {
				if !(bytes.Compare(sortKey, start) > 0 && bytes.Compare(sortKey, stop) < 0) {
					continue
				}
			}
			if c.noValue {
				result = append(result, []string{string(hashKey), string(sortKey)})
			} else {
				result = append(result, []string{string(hashKey), string(sortKey), string(value)})
			}
			keys[string(hashKey)] = append(keys[string(hashKey)], sortKey)
		}
	}
	if c.delete {
		for hk, sks := range keys {
			if err := ctx.table.MultiDel(context.Background(), []byte(hk), sks); err != nil {
				return err
			}
		}
	}
	if c.count {
		fmt.Println(len(result))
		return nil
	}
	table := tablewriter.NewWriter(os.Stdout)
	if c.noValue {
		table.SetHeader([]string{"hashKey", "sortKey"})
	} else {
		table.SetHeader([]string{"hashKey", "sortKey", "value"})
	}
	table.AppendBulk(result)
	table.Render()
	return nil
}

func (c *fullScanCommand) astNode() *ast.CommandASTNode {
	return &ast.CommandASTNode{
		SubNodes: map[string]*ast.CommandASTNode{},
	}
}
