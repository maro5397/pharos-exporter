package internal

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"syscall"
	"time"
)

type LogTailerConfig struct {
	Path         string
	PollInterval time.Duration
	Output       io.Writer
	FromStart    bool
}

type LogTailer struct {
	cfg    LogTailerConfig
	file   *os.File
	reader *bufio.Reader
	inode  uint64
	offset int64
}

func NewLogTailer(cfg LogTailerConfig) (*LogTailer, error) {
	if cfg.Path == "" {
		return nil, fmt.Errorf("log path is required")
	}
	if cfg.PollInterval <= 0 {
		cfg.PollInterval = time.Second
	}
	if cfg.Output == nil {
		cfg.Output = os.Stdout
	}
	return &LogTailer{cfg: cfg}, nil
}

func (t *LogTailer) Start(ctx context.Context) error {
	startAtEnd := !t.cfg.FromStart
	for {
		if err := t.openFile(startAtEnd); err != nil {
			if os.IsNotExist(err) {
				if err := sleepWithContext(ctx, t.cfg.PollInterval); err != nil {
					return err
				}
				continue
			}
			return err
		}
		break
	}
	defer t.closeFile()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line, err := t.reader.ReadBytes('\n')
		if len(line) > 0 {
			if _, werr := t.cfg.Output.Write(line); werr != nil {
				return werr
			}
			t.offset += int64(len(line))
		}
		if err == nil {
			continue
		}
		if err != io.EOF {
			return err
		}

		rotated, rerr := t.reopenIfRotated()
		if rerr != nil {
			return rerr
		}
		if rotated {
			continue
		}
		if err := sleepWithContext(ctx, t.cfg.PollInterval); err != nil {
			return err
		}
	}
}

func (t *LogTailer) reopenIfRotated() (bool, error) {
	info, err := os.Stat(t.cfg.Path)
	if err != nil {
		return false, err
	}
	inode, err := inodeFromInfo(info)
	if err != nil {
		return false, err
	}
	if inode != t.inode || info.Size() < t.offset {
		t.closeFile()
		if err := t.openFile(false); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func (t *LogTailer) openFile(startAtEnd bool) error {
	f, err := os.Open(t.cfg.Path)
	if err != nil {
		return err
	}
	info, err := f.Stat()
	if err != nil {
		f.Close()
		return err
	}
	inode, err := inodeFromInfo(info)
	if err != nil {
		f.Close()
		return err
	}
	offset := int64(0)
	if startAtEnd {
		if off, err := f.Seek(0, io.SeekEnd); err == nil {
			offset = off
		}
	}
	t.file = f
	t.reader = bufio.NewReader(f)
	t.inode = inode
	t.offset = offset
	return nil
}

func (t *LogTailer) closeFile() {
	if t.file != nil {
		_ = t.file.Close()
	}
	t.file = nil
	t.reader = nil
}

func inodeFromInfo(info os.FileInfo) (uint64, error) {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, fmt.Errorf("failed to read inode info")
	}
	return stat.Ino, nil
}
