package main

import (
	"fmt"
	"slices"
	"strings"
	"sync"
)

func startWorker(in, out chan interface{}, waiter *sync.WaitGroup, job job) {
	defer waiter.Done()
	defer close(out)
	job(in, out)
}

func ExecutePipeline(jobs ...job) {
	chan_len := len(jobs) + 1
	sl_ch := make([]chan interface{}, chan_len)
	for i := 0; i < chan_len; i++ {
		sl_ch[i] = make(chan interface{})
	}

	wg := &sync.WaitGroup{}
	for i, job := range jobs {
		wg.Add(1)
		go startWorker(sl_ch[i], sl_ch[i+1], wg, job)
	}
	// Wait нужен чтобы не выйти до результат выполнения функции а то я выйду, а тест не до конца сработает
	wg.Wait()

}

// func calc_crc32(in interface{}) string {
// 	data := fmt.Sprintf("%v", in)
// 	byte_data := []byte(data)
// 	hash := crc32.ChecksumIEEE([]byte(byte_data))
// 	return fmt.Sprintf("%v", hash)
// }
// func calc_md5(in interface{}) string {
// 	data := fmt.Sprintf("%v", in)
// 	hash := md5.Sum([]byte(data))
// 	return hex.EncodeToString(hash[:])

// }

func SingleHash(in, out chan interface{}) {
	// Косяк был что waitgroup тут не втюхал, и поэтому видиму функция завершалась до того как горутины отработают
	wg_main := &sync.WaitGroup{}
	for val := range in {
		str_val := fmt.Sprintf("%v", val)
		md5_data := DataSignerMd5(str_val)
		wg_main.Add(1)
		go func() {
			defer wg_main.Done()
			tmp_ch1 := make(chan interface{})
			tmp_ch2 := make(chan interface{})
			// fmt.Println("---------------")
			// fmt.Printf("%v-tmp_ch1: %v\n", str_val, &tmp_ch1)
			// fmt.Printf("%v-tmp_ch2: %v\n", str_val, &tmp_ch2)
			wg := &sync.WaitGroup{}
			wg.Add(2)
			go func() {
				defer wg.Done()
				tmp_ch1 <- DataSignerCrc32(md5_data)
			}()
			go func() {
				defer wg.Done()
				tmp_ch2 <- DataSignerCrc32(str_val)
			}()
			concat := fmt.Sprintf("%v~%v", <-tmp_ch2, <-tmp_ch1)
			wg.Wait() // По сути не нужен ведь я вверху пока не получу данные не пойду дальше
			out <- concat
		}()

	}
	wg_main.Wait()
}
func MultiHash(in, out chan interface{}) {
	wg_main := &sync.WaitGroup{}

	for val := range in {
		wg_main.Add(1)
		go func(tmp_val interface{}) {
			defer wg_main.Done()

			wg := &sync.WaitGroup{}
			result_slice := make([]string, 6)

			for i := 0; i <= 5; i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					concat_val := fmt.Sprintf("%v%v", i, tmp_val)
					crc32_val := DataSignerCrc32(concat_val)
					result_slice[i] = crc32_val

				}(i)
			}
			wg.Wait()
			// result_string = fmt.Sprintf("%v%v", result_string, crc32_val)
			result_string := strings.Join(result_slice, "")
			out <- result_string
		}(val)

	}
	wg_main.Wait()
}
func CombineResults(in, out chan interface{}) {
	result_slice := make([]string, 0, 7)
	for val := range in {
		result_slice = append(result_slice, fmt.Sprintf("%v", val))
	}
	slices.Sort(result_slice)
	result_string := strings.Join(result_slice, "_")
	out <- result_string
}

// Рабочий вариант, но без хитроебов с задержками и без использования их функций шифрования
// func SingleHash(in, out chan interface{}) {
// 	for val := range in {
// 		crc32_data := calc_crc32(val)
// 		crc32_md5_data := calc_crc32(calc_md5(val))
// 		concat := fmt.Sprintf("%v~%v", crc32_data, crc32_md5_data)
// 		out <- concat
// 	}

// }
// func MultiHash(in, out chan interface{}) {
// 	var result_string string
// 	for val := range in {
// 		result_string = ""
// 		for i := 0; i <= 5; i++ {
// 			concat_val := fmt.Sprintf("%v%v", i, val)
// 			crc32_val := calc_crc32(concat_val)
// 			result_string = fmt.Sprintf("%v%v", result_string, crc32_val)
// 		}
// 		out <- result_string
// 	}
// }
// func CombineResults(in, out chan interface{}) {
// 	result_slice := make([]string, 0, 7)
// 	for val := range in {
// 		result_slice = append(result_slice, fmt.Sprintf("%v", val))
// 	}
// 	slices.Sort(result_slice)
// 	result_string := strings.Join(result_slice, "_")
// 	out <- result_string
// }
