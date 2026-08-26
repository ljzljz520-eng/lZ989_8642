package service

import (
	"context"
	"customerfollowup/internal/model"
	"sort"
	"strings"
)

type Report struct {
	Records  []model.Record
	ByStatus map[string]int
	Tags     map[string]int
}

func (x *Service) Report(ctx context.Context, customer string) (Report, error) {
	rows, e := x.Query(ctx, customer)
	if e != nil {
		return Report{}, e
	}
	r := Report{Records: rows, ByStatus: map[string]int{}, Tags: map[string]int{}}
	for _, v := range rows {
		r.ByStatus[v.Status]++
		for _, tag := range model.NormalizeTags(v.Tags) {
			r.Tags[tag]++
		}
	}
	sort.Slice(r.Records, func(i, j int) bool { return r.Records[i].UpdatedAt.Before(r.Records[j].UpdatedAt) })
	return r, nil
}
func FilterStatus(rows []model.Record, status string) []model.Record {
	out := []model.Record{}
	for _, r := range rows {
		if r.Status == status {
			out = append(out, r)
		}
	}
	return out
}
func FilterTag(rows []model.Record, tag string) []model.Record {
	tag = strings.ToLower(strings.TrimSpace(tag))
	out := []model.Record{}
	for _, r := range rows {
		for _, v := range model.NormalizeTags(r.Tags) {
			if v == tag {
				out = append(out, r)
				break
			}
		}
	}
	return out
}
func MergeReports(a, b Report) Report {
	out := Report{Records: append(append([]model.Record{}, a.Records...), b.Records...), ByStatus: map[string]int{}, Tags: map[string]int{}}
	for k, v := range a.ByStatus {
		out.ByStatus[k] += v
	}
	for k, v := range b.ByStatus {
		out.ByStatus[k] += v
	}
	for k, v := range a.Tags {
		out.Tags[k] += v
	}
	for k, v := range b.Tags {
		out.Tags[k] += v
	}
	return out
}
func Statuses() []string { return []string{"", "processed", "archived"} }
func IsKnownStatus(s string) bool {
	for _, v := range Statuses() {
		if s == v {
			return true
		}
	}
	return false
}
func NormalizeRecord(r model.Record) model.Record {
	r.Tags = model.NormalizeTags(r.Tags)
	r.Summary = strings.TrimSpace(r.Summary)
	return r
}
func (x *Service) NormalizeAndSave(ctx context.Context, r model.Record) error {
	return x.Register(ctx, NormalizeRecord(r))
}
func (x *Service) ArchiveMany(ctx context.Context, ids []string, actor string) int {
	n := 0
	for _, id := range ids {
		if x.Archive(ctx, id, actor) == nil {
			n++
		}
	}
	return n
}
func (x *Service) ProcessMany(ctx context.Context, ids []string, actor string) int {
	n := 0
	for _, id := range ids {
		if x.Process(ctx, id, actor) == nil {
			n++
		}
	}
	return n
}
func (x *Service) FindByTag(ctx context.Context, tag string) ([]model.Record, error) {
	r, e := x.Report(ctx, "")
	if e != nil {
		return nil, e
	}
	return FilterTag(r.Records, tag), nil
}
func (x *Service) FindByStatus(ctx context.Context, status string) ([]model.Record, error) {
	r, e := x.Report(ctx, "")
	if e != nil {
		return nil, e
	}
	return FilterStatus(r.Records, status), nil
}
