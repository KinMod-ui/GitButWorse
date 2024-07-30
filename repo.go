package main


type repo struct {
    basePath string
}

func (r repo) get_root() string{
   return r.basePath 
}
