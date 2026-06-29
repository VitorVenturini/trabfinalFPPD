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