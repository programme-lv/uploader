GEN = ENV["GEN_EXE"]
SOL = ENV["SOL_EXE"]

N_LIM = 100000
TESTS = 5

for i = 1:TESTS
    # testname 0001, 0002 and so on
    testname = string(i,base=10,pad=3)
    N = ((i/TESTS)^2)*N_LIM
    open(f->write(f, read(`$GEN $(N)`, String)), "$(testname).in", "w")
    open(f->write(f, read(pipeline(`cat $(testname).in`,`$SOL`), String)), "$(testname).ans", "w")
end

