package model

import (
	"fmt"
)

func GetVote(id int64) *Vote {
	ret := &Vote{}
	if err := GlobalConn.Preload("VoteOpt").Table("vote_title").Where("id = ?", id).First(ret).Error; err != nil {
		fmt.Printf("err:%s\n", err.Error())
	}
	return ret
}
func GetVotes() []*Vote {
	ret := make([]*Vote, 0)
	sql := "select * from `vote_title` where id>0"
	err := GlobalConn.Raw(sql).Scan(&ret).Error
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
		return ret
	}
	for _, vote := range ret {
		vote.VoteOpt = make([]*VoteOpt, 0)
	}
	return ret
}
func CreateVote(vote *Vote) error {
	sql := "insert into `vote_title` values(?,?,?,?,?,?,?,?)"
	return GlobalConn.Exec(sql, vote.Id, vote.Title, vote.StartTime, vote.Type, vote.UserId, vote.Count, vote.CreatedTime, vote.During).Error
}
func UpdateVote(vote *Vote) error {
	sql := "update `vote_title` set title=?,start_time=?,type=?,user_id=?,count=?,created_time=?,during=? where id=? "
	return GlobalConn.Exec(sql, vote.Title, vote.StartTime, vote.Type, vote.UserId, vote.Count, vote.CreatedTime, vote.During, vote.Id).Error
}
func DeleteVote(id int64) error {
	sql := "delete from `vote_title` where id=? "
	return GlobalConn.Exec(sql, id).Error
}

func DoVote(userId, voteId int64, opt []int64) bool {
	var err error
	tx := GlobalConn.Begin()
	var count int64
	sql := "select count(id) from vote_opt_user where user_id = ? and vote_id = ?"
	err = tx.Raw(sql, userId, voteId).Scan(&count).Error
	if err != nil || count > 0 {
		//fmt.Printf("err:%s", err.Error())
		tx.Rollback()
		return false
	}
	sql2 := "update `vote_title` set count=count+? where id=?"
	err = tx.Exec(sql2, len(opt), voteId).Error
	if err != nil {
		fmt.Printf("err:%s", err.Error())
		tx.Rollback()
		return false
	}

	for _, va := range opt {
		sql3 := "update `vote_opt` set count=count+1 where id=?"
		err = tx.Exec(sql3, va).Error
		user := VoteOptUser{
			VoteId:    voteId,
			UserId:    userId,
			VoteOptId: va,
		}
		err = tx.Create(&user).Error
	}
	if err != nil {
		fmt.Printf("err:%s", err.Error())
		tx.Rollback()
		return false
	}
	tx.Commit()
	return true
}
func GetOpts(voteId int64) []*VoteOpt {
	ret := make([]*VoteOpt, 0, 0)
	sql := "select * from `vote_opt` where vote_id=?"
	err := GlobalConn.Raw(sql, voteId).Scan(&ret).Error
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
	}
	return ret
}
func CreateOpt(opt *VoteOpt) error {
	sql := "insert into `vote_opt` values (?, ? ,? ,? ,? ,? )"
	return GlobalConn.Exec(sql, opt.Id, opt.VoteId, opt.Count, opt.Name, opt.CreatedTime, opt.Percent).Error
}
func UpdateOpt(opt *VoteOpt) error {
	sql := "update `vote_opt` set vote_id=?,count=?,name=?,created_time=?,percent=? where id=? "
	return GlobalConn.Exec(sql, opt.VoteId, opt.Count, opt.Name, opt.CreatedTime, opt.Percent, opt.Id).Error
}
func DeleteOpt(id int64) error {
	sql := "delete from `vote_opt` where id=?"
	return GlobalConn.Exec(sql, id).Error
}
