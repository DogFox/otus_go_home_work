package hw06_pipeline_execution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func worker(stageResultCh In, done In) Bi {
	// stageResultCh <- value мы так делать не можем, потому что канал на чтение, но можно переопределить новым каналом
	// для чтения и записи, но сигнатура при этом не поменяется
	newResultCh := make(Bi)
	go func() {
		defer close(newResultCh)
		// подождали результат, завернули в новый канал, вернули его
		for {
			select {
			case value, ok := <-stageResultCh:
				if !ok {
					return
				}
				newResultCh <- value
			case <-done:
				return
			}
		}
	}()
	return newResultCh
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// перебираем стейджи
	for _, stage := range stages {
		// отдаем входной канал с данными и переопределяем его же после выполнения стейджа
		// на следующей итерации в нем лежит результат выполнения предыдущей итерации
		in = worker(stage(in), done)
	}
	return in
}
