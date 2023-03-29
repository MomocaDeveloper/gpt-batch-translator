package services

import (
	"strings"
	"errors"
	"sync"
)

type WorkerPool struct{
	Id string `json:"id"`
	Total int `json:"total"`
	Finish int `json:"finish"`
	MaxWorkers int //最大的goroutine数量
	WorkerChan chan struct{}   //控制goroutine数量的通道
	WaitGroup sync.WaitGroup
}

type Progress struct{
	AllPools map[string]*WorkerPool `json:"allPools"`
}

func removeSuffix(fileName string) string {
	parts := strings.Split(fileName, ".")
	if len(parts) > 1 {
		return strings.Join(parts[:len(parts)-1], ".")
	} else {
		return fileName
	}
}

func (p *Progress)CreateMsgData()map[string]int{
	if len(p.AllPools) == 0{
		return map[string]int{}
	}
	msg := make(map[string]int)
	for _, pool := range p.AllPools{
		fid := removeSuffix(pool.Id)
		msg[fid] = pool.Finish*100/pool.Total
	}
	return msg
}

func (p *Progress)UpdateProgress(id string, count int){
	if pool, ok := p.AllPools[id];ok {
		pool.Finish += count
		p.AllPools[id] = pool
	}
}

func (p *Progress)CreateNewPool(id string, maxWorkers, totalLine int)error{
	if _, ok := p.AllPools[id];ok{
		return errors.New("the pool already exist!")
	}
	newPool := &WorkerPool{
		Id: id,
		MaxWorkers:maxWorkers,
		Total : totalLine,
		Finish: 0,
		WorkerChan:make(chan struct{}, maxWorkers),
	}
	p.AllPools[id] = newPool
	return nil
}

func (p *Progress)FinishPool(id string){
	delete(p.AllPools, id)
}

func (p *Progress)AlreadyFinishPool(id string)bool{
	pool, ok := p.AllPools[id]
	if !ok{
		return true
	}
	return pool.Total<=pool.Finish
}

func (p *Progress)GetFinishCount(id string)int{
	pool, ok := p.AllPools[id]
	if !ok{
		return 0
	}
	return pool.Finish
}

func (p *Progress)GetWorkerPool(id string)(*WorkerPool, bool){
	pool, err := p.AllPools[id]
	return pool, err
}

func (p *WorkerPool)AddTask(task func(param *Content), param *Content){
	p.WaitGroup.Add(1)
	p.WorkerChan <- struct{}{}
	go func(){
		task(param)
		<-p.WorkerChan
		p.WaitGroup.Done()
	}()
}

func (p *WorkerPool)Wait(){
	p.WaitGroup.Wait()
}