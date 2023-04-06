%-
operators: 
    -eq         :: int -> bool
    -eqstr      :: string -> bool
    -lt         :: int -> bool
    -lte        :: int -> bool
    -gt         :: int -> bool
    -gte        :: int -> bool
    -not        :: bool -> bool
-% 

echo "Hello world"

if -eq 5 5 then

end

if -not -eq 5 5 then 

end

if -not ( -eq 5 5 ) then

end

proc A do
    
    echo $1 $2

end
