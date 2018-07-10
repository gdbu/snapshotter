package main

import (
	"flag"
	"fmt"
	"path"
	"time"

	"github.com/Hatch1fy/pgutils"
	"github.com/Hatch1fy/snapshotter/frontends"

	"github.com/Hatch1fy/snapshotter"
	"github.com/Hatch1fy/snapshotter/backends"
	"github.com/missionMeteora/journaler"
	"github.com/missionMeteora/toolkit/closer"
)

func main() {
	var (
		s  *snapshotter.Snapshotter
		fe snapshotter.Frontend
		be snapshotter.Backend

		cfgPath string

		cfg   Config
		pgcfg pgutils.Config
		s3cfg backends.S3Config
		sscfg snapshotter.Config

		err error
	)

	flag.StringVar(&cfgPath, "config", "./cfg", "Path of configuration files")
	flag.Parse()

	out := journaler.New("Postgres snapshotter")
	out.Notification("Starting service, Hello sir!")

	if cfg, err = newConfig(path.Join(cfgPath, "config.toml")); err != nil {
		out.Error("Error parsing configuration: %v", err)
		return
	}

	if pgcfg, err = pgutils.NewConfig(path.Join(cfgPath, "postgres.toml")); err != nil {
		out.Error("Error parsing Postgres configuration: %v", err)
		return
	}

	if s3cfg, err = backends.NewS3Config(path.Join(cfgPath, "s3.toml")); err != nil {
		out.Error("Error parsing S3 configuration: %v", err)
		return
	}

	sscfg.Extension = "sql"
	sscfg.Name = cfg.Name
	sscfg.Interval = cfg.Interval * time.Minute
	sscfg.Truncate = time.Hour

	fe = frontends.NewPostgres(pgcfg)

	s3bucket := fmt.Sprintf("%s.%s", cfg.Bucket, cfg.Environment)

	if be, err = backends.NewS3(s3cfg.Config(), s3bucket); err != nil {
		out.Error("Error creating S3 backend: %v", err)
		return
	}

	if s, err = snapshotter.New(fe, be, sscfg); err != nil {
		out.Error("Error starting snapshotter:", err)
		return
	}

	c := closer.New()
	c.Wait()
	journaler.Notification("Closing service, see you again soon!")
	s.Close()
}
