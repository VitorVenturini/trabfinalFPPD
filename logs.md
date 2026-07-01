fppd3217@atlantica:~/trabfinalFPPD$ go run ./sequencial -n 4 -seed 42
Running sequential matrix multiplication with n=4
mode=sequential n=4 seed=42 nodes=1 processes=1 elapsed_sec=0.000001 label=""
verification:
  c00=0.761944217240602
  c0n=0.606335030831764
  cn0=0.915571669715378
  cnn=0.801411395076022
  checksum=12.234046811807353


fppd3217@atlantica:~/trabfinalFPPD$ go run ./sequencial -n 8 -seed 42
Running sequential matrix multiplication with n=8
mode=sequential n=8 seed=42 nodes=1 processes=1 elapsed_sec=0.000002 label=""
verification:
  c00=1.210975612040805
  c0n=1.629680338003836
  cn0=1.337132415594078
  cnn=1.736407089868683
  checksum=133.614559094271414

fppd3217@atlantica:~/trabfinalFPPD$ go run ./sequencial -n 16 -seed 42
Running sequential matrix multiplication with n=16
mode=sequential n=16 seed=42 nodes=1 processes=1 elapsed_sec=0.000024 label=""
verification:
  c00=2.382912377393811
  c0n=3.619603125102645
  cn0=4.454654765527576
  cnn=5.114767047032694
  checksum=1007.669016598953249

fppd3217@atlantica:~/trabfinalFPPD$ go run ./sequencial -n 1024 -seed 42
Running sequential matrix multiplication with n=1024
mode=sequential n=1024 seed=42 nodes=1 processes=1 elapsed_sec=4.475530 label=""
verification:
  c00=261.986299041784548
  c0n=262.025978125949848
  cn0=257.204914579881972
  cnn=259.657684276494763
  checksum=268678649.415891289710999


PARALELO

fppd3217@atlantica:~/trabfinalFPPD$ mpirun -np 1 go run ./paralelo -n 4 -seed 42
Running parallel matrix multiplication with n=4 and processes=1
mode=parallel n=4 seed=42 nodes=1 processes=1 elapsed_sec=0.000002 label=""
verification:
  c00=0.761944217240602
  c0n=0.606335030831764
  cn0=0.915571669715378
  cnn=0.801411395076022
  checksum=12.234046811807353

fppd3217@atlantica:~/trabfinalFPPD$ mpirun -np 2 go run ./paralelo -n 4 -seed 42
Running parallel matrix multiplication with n=4 and processes=2
mode=parallel n=4 seed=42 nodes=1 processes=2 elapsed_sec=0.000292 label=""
verification:
  c00=0.761944217240602
  c0n=0.606335030831764
  cn0=0.915571669715378
  cnn=0.801411395076022
  checksum=12.234046811807353


DEU ERRO DE FALTA DE SLOTS E PRECISOU DO OVERSUBSCRIBE
fppd3217@atlantica:~/trabfinalFPPD$ mpirun --oversubscribe -np 3 go run ./paralelo -n 8 -seed 42
Running parallel matrix multiplication with n=8 and processes=3
mode=parallel n=8 seed=42 nodes=1 processes=3 elapsed_sec=0.004423 label=""
verification:
  c00=1.210975612040805
  c0n=1.629680338003836
  cn0=1.337132415594078
  cnn=1.736407089868683
  checksum=133.614559094271414


fppd3217@atlantica:~/trabfinalFPPD$ mpirun --oversubscribe -np 4 go run ./paralelo -n 16 -seed 42
Running parallel matrix multiplication with n=16 and processes=4
mode=parallel n=16 seed=42 nodes=1 processes=4 elapsed_sec=0.003651 label=""
verification:
  c00=2.382912377393811
  c0n=3.619603125102645
  cn0=4.454654765527576
  cnn=5.114767047032694
  checksum=1007.669016598953249