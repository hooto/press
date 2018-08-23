
{{range $v := .config.Buckets}}
source {{$v.Name}}_full
{
    type                = xmlpipe2
    xmlpipe_fixup_utf8  = 1
    xmlpipe_command     = cat {{$.config.Prefix}}/var/sphinxsearch/{{$v.Name}}/full.xml
}

source {{$v.Name}}_delta : {{$v.Name}}_full
{
    xmlpipe_command     = cat {{$.config.Prefix}}/var/sphinxsearch/{{$v.Name}}/delta.xml
}

index {{$v.Name}}_full
{
    source          = {{$v.Name}}_full

    path            = {{$.config.Prefix}}/var/sphinxsearch/{{$v.Name}}/full
    docinfo         = extern
    
    mlock           = 0
    morphology      = none

    stopwords       = {{$.config.Prefix}}/misc/sphinxsearch/stopword.conf
    min_word_len    = 1
    charset_table   = 0..9, A..Z->a..z, _, a..z, U+410..U+42F->U+430..U+44F, U+430..U+44F
    ngram_len       = 1
    ngram_chars     = U+3000..U+2FA1F

    phrase_boundary_step    = 100
    html_strip              = 0
}

index {{$v.Name}}_delta : {{$v.Name}}_full
{
    source          = {{$v.Name}}_delta
    path            = {{$.config.Prefix}}/var/sphinxsearch/{{$v.Name}}/delta
}

index {{$v.Name}}
{
    type  = distributed
    local = {{$v.Name}}_full
    local = {{$v.Name}}_delta
}

{{end}}


indexer
{
    mem_limit             = 128M
    max_xmlpipe2_field    = 2M
    write_buffer          = 8M
    max_iops              = 40
    max_iosize            = 1048576
    max_file_field_buffer = 64M
}

searchd
{
    listen           = {{.config.Prefix}}/var/sphinxsearch/searchd.sock
    log              = {{.config.Prefix}}/var/sphinxsearch/searchd.log
    read_timeout     = 5
    client_timeout   = 300    
    pid_file         = {{.config.Prefix}}/var/sphinxsearch/searchd.pid
    read_buffer      = 4M
    mva_updates_pool = 16M
    seamless_rotate  = 1
    preopen_indexes  = 1    
    unlink_old       = 1    
    workers          = prefork 
    max_packet_size  = 16M
}

