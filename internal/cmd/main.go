// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {
	// setupLog := ctrl.Log.WithName("experimental-launcher")

	// Graceful shutdown context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start experimental controllers
	// go func() {
	// 	setupLog.Info("Starting experimental controllers")
	// 	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
	// 		Scheme: runtime.Scheme,
	// 	})
	// 	if err != nil {
	// 		setupLog.Error(err, "unable to start manager for experimental controllers")
	// 		cancel()
	// 		os.Exit(1)
	// 	}

	// 	// Register experimental controllers
	// 	if err := (&experimental.ExperimentalController{}).SetupWithManager(mgr); err != nil {
	// 		setupLog.Error(err, "unable to setup experimental controller")
	// 		cancel()
	// 		os.Exit(1)
	// 	}

	// 	// Start the experimental manager
	// 	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
	// 		setupLog.Error(err, "problem running experimental manager")
	// 		cancel()
	// 		os.Exit(1)
	// 	}
	// }()

	// Launch production binary as a subprocess
	cmd := exec.Command("./my-operator-production-binary") // Path to your production binary
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the subprocess
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start production binary: %v", err)
		cancel()
		os.Exit(1)
	}

	log.Println("Experimental launcher is running. Production binary started.")

	// Handle OS signals for graceful shutdown
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	select {
	case <-signalCh: // If the launcher gets terminated
		log.Println("Received termination signal. Stopping everything...")

		// Stop the experimental controllers by calling cancel()
		cancel()

		// Kill the production subprocess
		if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
			log.Printf("Failed to terminate production subprocess: %v", err)
		}

	case <-ctx.Done(): // If the experimental controllers fail
		log.Println("Experimental controllers stopped. Terminating...")
		if err := cmd.Process.Kill(); err != nil {
			log.Printf("Failed to force-kill production subprocess: %v", err)
		}
	}

	// Wait for subprocess to fully shut down
	_ = cmd.Wait()

	log.Println("Experimental launcher shut down completely.")
}

// func startExperimentalControllers(setupLog logr.Logger, cancel context.CancelCauseFunc) error {
// 	setupLog.Info("Starting experimental controllers")
// 	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
// 		Scheme: ctrl.Scheme,
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to start experimentar controller manager: %w", err)
// 	}

// 	// TODO Register experimental controllers
// 	if err := (&experimental.ExperimentalController{}).SetupWithManager(mgr); err != nil {
// 		return fmt.Errorf("failed to setup experimental controller: %w", err)
// 	}

// 	go func() {
// 		if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
// 			setupLog.Error(err, "problem running experimental manager")
// 			cancel()
// 			os.Exit(1)
// 		}
// 	}()
// 	return nil
// }
